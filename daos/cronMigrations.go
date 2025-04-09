package daos

import (
	"CollegeAdministration/models"
	"log"
)

func (AC *Daos) RunMigrationsForRemovingOutDatedTokens() error {

	err := AC.dbConn.Where("is_valid = ?", false).Delete(&models.Token_generator{}).Error
	if err != nil {
		log.Println("Error run cron migrations", err)
		return err
	}
	return nil
}

func (AC *Daos) GetAllTokens() ([]*models.Token_generator, error) {

	var all_tokens []*models.Token_generator
	err := AC.dbConn.Select("*").Find(&all_tokens).Error
	if err != nil {
		return nil, err
	}

	return all_tokens, nil
}
