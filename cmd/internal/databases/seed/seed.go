package seed

import "gorm.io/gorm"

func AllSeed(db *gorm.DB) error {

	err := SystemSettingSeed(db)

	if err != nil {
		return err
	}

	err = PrivacyPolicySeed(db)

	if err != nil {
		return err
	}

	err = UserRoleSeed(db)

	if err != nil {
		return err
	}

	err = RightGroupSeed(db)

	if err != nil {
		return err
	}

	err = RightCategorySeed(db)

	if err != nil {
		return err
	}

	err = RightDetailSeed(db)

	if err != nil {
		return err
	}

	err = RegionSeed(db)

	if err != nil {
		return err
	}

	err = SubRegionSeed(db)

	if err != nil {
		return err
	}

	err = CountrySeed(db)

	if err != nil {
		return err
	}

	err = StateSeed(db)

	if err != nil {
		return err
	}

	err = CitySeed(db)

	if err != nil {
		return err
	}

	err = ActivityList(db)

	if err != nil {
		return err
	}

	err = EmailSettingSeed(db)

	if err != nil {
		return err
	}

	err = BackupType(db)

	if err != nil {
		return err
	}

	err = AuthenticationSettingSeed(db)

	if err != nil {
		return err
	}

	err = SmsSettingSeed(db)

	if err != nil {
		return err
	}

	return nil
}
