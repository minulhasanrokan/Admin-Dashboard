package seed

import (
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/models"

	"gorm.io/gorm"
)

func PrivacyPolicySeed(db *gorm.DB) error {

	var count int64

	db.Model(&models.PrivacyPolicy{}).Count(&count)

	if count > 0 {
		return nil
	}

	data := models.PrivacyPolicy{

		MainBaseModel: BaseSeedData(),
		Content:       helpers.Str("\n\nAt Smart HR, we value your privacy and are committed to protecting your personal data. This Privacy Policy explains how we collect, use, disclose, and safeguard your information when you access our HR admin platform. Please read this policy carefully.\nInformation We Collect\n\nWe may collect the following types of personal information:\n\n    Employee Data :Name, contact information, job title, department, payroll details, and other employment-related data.\n    Login Information :Username, password, IP address, and other access credentials.\n    Usage Data : Details about how the platform is used, including logs of actions taken within the HR admin module\n\nHow We Use Your Information\n\nWe may use the information we collect to:\n\n    Manage employee records and HR functions.\n    Process payroll and benefits administration.\n    Ensure the security and integrity of the platform.\n    Generate reports and analytics to improve our services.\n\nSharing Your Information\n\nWe will not sell or share your personal information with third parties except:\n\n    To authorized service providers who assist in our HR processes.\n    If required by law, such as in response to a court order or legal obligation.\n\nData Security\n\nWe implement appropriate technical and organizational measures to protect your personal data from unauthorized access, alteration, or disclosure. However, no method of transmission over the internet is 100% secure, and we cannot guarantee absolute security.\nData Retention\n\nWe will retain your personal information for as long as necessary to fulfill the purposes outlined in this policy, or as required by law.\nYour Rights\n\nYou may have the following rights concerning your personal data:\n\n    Access to your data.\n    Correction of inaccurate information.\n    Deletion of data where appropriate.\n    Restriction of processing under certain circumstances.\n\nChanges to This Policy\n\nWe may update this policy from time to time. The latest version will always be available on the platform. Continued use of the HR admin platform signifies your acceptance of any updates.\n"),
	}

	if err := db.Create(&data).Error; err != nil {
		return err
	}

	history := models.PrivacyPolicyHistory{
		PrivacyPolicy: data,
		MainId:        data.Id,
	}

	return db.Create(&history).Error
}
