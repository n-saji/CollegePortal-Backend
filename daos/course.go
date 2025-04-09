package daos

import (
	"CollegeAdministration/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func (ac *Daos) InsertValuesToCoursesAvailable(ca *models.CourseInfo) error {

	err := ac.dbConn.Table("course_infos").Create(ca).Error
	if err != nil {
		log.Println("not able to insert to course_infos table ", err)
		return fmt.Errorf("failed! %s", err.Error())
	}
	return nil

}
func (ac *Daos) GetCourseByName(name string) (models.CourseInfo, error) {

	var ca models.CourseInfo
	val := ac.dbConn.Select("*").Table("course_infos").Where("course_name = ?", name).First(&ca)
	if val.Error != nil {
		log.Println("Not able to Fetch value from  table ", val.Error)
		return ca, val.Error
	}
	return ca, nil
}

func (ac *Daos) GetCourseById(id uuid.UUID) (models.CourseInfo, error) {

	var ca models.CourseInfo
	// if id == uuid.Nil {
	// 	return ca, fmt.Errorf("UUID is NULL for course ID! Add new Course to ")
	// }
	val := ac.dbConn.Select("*").Table("course_infos").Where("id = ?", id).First(&ca)
	if val.Error != nil {
		log.Println("Not able to select from course_infos table ", val.Error)
		return ca, val.Error
	}
	return ca, nil
}

func (ac *Daos) RetieveCoursesAvailable() ([]*models.CourseInfo, error) {

	var rca []*models.CourseInfo
	err := ac.dbConn.Order("course_name").Find(&rca).Error

	return rca, err

}
func (ac *Daos) UpdateCourseByName(name string, rc *models.CourseInfo) error {
	rc_existing, _ := ac.GetCourseByName(rc.CourseName)
	if rc_existing.Id != uuid.Nil {
		return fmt.Errorf("course already exits")
	}
	rcOld, err1 := ac.GetCourseByName(name)
	if err1 != nil {
		return fmt.Errorf("course not available %s", err1.Error())
	}
	rc.Id = rcOld.Id
	err := ac.dbConn.Save(&rc).Error
	if err != nil {
		return fmt.Errorf("failed UpdateCourseByName %s", err.Error())
	}
	return nil
}

func (ac *Daos) CheckCourse(coursename string) bool {

	var len int64
	ac.dbConn.Model(models.CourseInfo{}).Where("course_name = ?", coursename).Count(&len)
	if len > 0 {
		return true
	} else {
		return false
	}
}

func (ac *Daos) DeleteCourse(id uuid.UUID) (bool, error) {

	err := ac.dbConn.Where("id = ?", id).Delete(&models.CourseInfo{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
