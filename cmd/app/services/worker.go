package services

import (
	"admin-dashboard/cmd/app/repositories"
	"context"
	"fmt"
	"time"
)

func (s *Service) ProcessPendingEmails(ctx context.Context, batchSize int, maxTries int) (int, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	pendingMails, err := repository.GetPendingMailHistories(batchSize, maxTries)
	if err != nil {
		return 0, err
	}

	processedCount := 0

	for _, mailItem := range pendingMails {
		select {
		case <-ctx.Done():
			return processedCount, ctx.Err()
		default:
		}

		mailItem.TryStatus++

		sendErr := s.SendEmailDirect(mailItem.ToAddress, mailItem.Subject, mailItem.Body)
		if sendErr != nil {
			mailItem.SentStatus = 0
			mailItem.ErrorStatus = 1
			mailItem.ErrorMsg = sendErr.Error()
			if s.Logger != nil {
				s.Logger.Warnf("[WORKER] Failed to send email ID %d to %s (attempt %d): %v", mailItem.Id, mailItem.ToAddress, mailItem.TryStatus, sendErr)
			}
		} else {
			mailItem.SentStatus = 1
			mailItem.ErrorStatus = 0
			mailItem.ErrorMsg = ""
			if s.Logger != nil {
				s.Logger.Infof("[WORKER] Successfully sent email ID %d to %s", mailItem.Id, mailItem.ToAddress)
			}
		}

		_ = repository.UpdateMailHistory(&mailItem)
		processedCount++
	}

	return processedCount, nil
}

func (s *Service) ProcessPendingSms(ctx context.Context, batchSize int, maxTries int) (int, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	pendingSms, err := repository.GetPendingSmsHistories(batchSize, maxTries)
	if err != nil {
		return 0, err
	}

	processedCount := 0

	for _, smsItem := range pendingSms {
		select {
		case <-ctx.Done():
			return processedCount, ctx.Err()
		default:
		}

		smsItem.TryStatus++

		_, sendErr := s.SendSmsDirect(smsItem.ToMobile, smsItem.Message)
		if sendErr != nil {
			smsItem.SentStatus = 0
			smsItem.ErrorStatus = 1
			smsItem.ErrorMsg = sendErr.Error()
			if s.Logger != nil {
				s.Logger.Warnf("[WORKER] Failed to send SMS ID %d to %s (attempt %d): %v", smsItem.Id, smsItem.ToMobile, smsItem.TryStatus, sendErr)
			}
		} else {
			smsItem.SentStatus = 1
			smsItem.ErrorStatus = 0
			smsItem.ErrorMsg = ""
			if s.Logger != nil {
				s.Logger.Infof("[WORKER] Successfully sent SMS ID %d to %s", smsItem.Id, smsItem.ToMobile)
			}
		}

		_ = repository.UpdateSmsHistory(&smsItem)
		processedCount++
	}

	return processedCount, nil
}

func (s *Service) ProcessPendingQueue(ctx context.Context) (int, int, error) {

	batchSize := s.Config.QueueWorkerBatchSize
	if batchSize <= 0 {
		batchSize = 50
	}

	maxTries := s.Config.QueueWorkerMaxTries
	if maxTries <= 0 {
		maxTries = 5
	}

	emailCount, errEmail := s.ProcessPendingEmails(ctx, batchSize, maxTries)
	smsCount, errSms := s.ProcessPendingSms(ctx, batchSize, maxTries)

	if errEmail != nil && errSms != nil {
		return emailCount, smsCount, fmt.Errorf("email err: %w; sms err: %v", errEmail, errSms)
	}
	if errEmail != nil {
		return emailCount, smsCount, errEmail
	}
	if errSms != nil {
		return emailCount, smsCount, errSms
	}

	return emailCount, smsCount, nil
}

func (s *Service) StartWorker(ctx context.Context, interval time.Duration) {

	if interval <= 0 {
		if s.Config.QueueWorkerInterval > 0 {
			interval = time.Duration(s.Config.QueueWorkerInterval) * time.Second
		} else {
			interval = 5 * time.Second
		}
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if s.Logger != nil {
		s.Logger.Infof("[WORKER] Pending queue worker started, polling every %v", interval)
	}

	_, _, _ = s.ProcessPendingQueue(ctx)

	for {
		select {
		case <-ctx.Done():
			if s.Logger != nil {
				s.Logger.Info("[WORKER] Pending queue worker stopping...")
			}
			return
		case <-ticker.C:
			_, _, _ = s.ProcessPendingQueue(ctx)
		}
	}
}
