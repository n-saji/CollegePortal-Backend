package service

import (
	"CollegeAdministration/config"
	"CollegeAdministration/models"
	"CollegeAdministration/utils"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
)

func (ac *Service) InsertValuesToCAd(new_student *models.StudentInfo, account_id string) error {

	course_details, err := ac.daos.GetCourseByName(new_student.ClassesEnrolled.CourseName)
	if err != nil {
		return fmt.Errorf("course Not Found")
	}
	if _, err := strconv.ParseFloat(new_student.Name, 64); err == nil {
		return fmt.Errorf("name cant be number")
	}
	new_student.ClassesEnrolled.Id = course_details.Id
	new_student.ClassesEnrolled = course_details
	new_student.CourseId = course_details.Id
	sd_existing, _ := ac.daos.GetStudentDetailsByRollNumber(new_student.RollNumber)

	for _, each_student := range sd_existing {

		course_details, _ := ac.daos.GetCourseById(each_student.CourseId)
		each_student.ClassesEnrolled = course_details

		if each_student.RollNumber == new_student.RollNumber {
			if each_student.Name != new_student.Name {
				return fmt.Errorf(fmt.Sprintf("student %s already present with roll number %s", each_student.Name, new_student.RollNumber))
			}
			if each_student.ClassesEnrolled.CourseName == new_student.ClassesEnrolled.CourseName {
				return fmt.Errorf(fmt.Sprintf("student already present with course %s", each_student.ClassesEnrolled.CourseName))
			}

			if each_student.Age != new_student.Age {
				return fmt.Errorf(fmt.Sprintf("student age  mismatch exisiting age %d", each_student.Age))
			}
		}

	}
	sd, _ := ac.daos.GetStudentdetail(new_student)
	if sd != nil && sd.CourseId == new_student.CourseId {
		return fmt.Errorf("student with course exist")
	}

	new_student.Id = uuid.New()
	sm, err2 := ac.InsertStudentIdInToMarksTable(new_student)
	if err2 != nil {
		return err2
	}
	new_student.MarksId = sm.Id

	err1 := ac.daos.InsertValuesToCollegeAdministration(new_student)
	if err1 != nil {
		return err1
	}

	instructor_list_old, _ := ac.GetInstructorDetailWithSpecifics(models.InstructorDetails{CourseId: new_student.CourseId})
	for _, each_instructor := range instructor_list_old {
		err := ac.Update_Instructor_Info(each_instructor, models.InstructorDetails{Id: each_instructor.Id})
		if err != nil {
			return err
		}
	}
	account := &models.Account{}
	account.Id = new_student.Id
	account.Info.Credentials.Id = new_student.Id
	account.Name = new_student.Name
	account.Type = config.AccountTypeStudent
	account.Verified = false

	err3 := ac.daos.CreateAccount(account)
	if err3 != nil {
		log.Println("error storing in account student")
		return err3
	}
	go utils.StoreMessages("New Student", account.Name, config.AccountTypeInstructor, account_id)
	return nil

}

func (ac *Service) Retrieve_student_details() ([]*models.StudentInfo, error) {

	rca, err := ac.daos.RetrieveCollegeAdministration()
	for _, each_student := range rca {
		if each_student.ClassesEnrolled.CourseName == "" {
			deleted_course, _ := ac.daos.GetCourseByName("Course Deleted")
			each_student.ClassesEnrolled = deleted_course
			each_student.CourseId = deleted_course.Id
			ac.daos.UpdateClgStudent(each_student)
		}
	}
	if err != nil {
		return rca, err
	}
	return rca, nil
}

func (ac *Service) Retrieve_student_detailsbyOrder(order string) ([]*models.StudentInfo, error) {

	rca, err := ac.daos.RetrieveCollegeAdministrationByOrder(order)
	for _, each_student := range rca {
		if each_student.ClassesEnrolled.CourseName == "" {
			deleted_course, _ := ac.daos.GetCourseByName("Course Deleted")
			each_student.ClassesEnrolled = deleted_course
			each_student.CourseId = deleted_course.Id
			ac.daos.UpdateClgStudent(each_student)
		}
	}
	if err != nil {
		return rca, err
	}
	return rca, nil
}

func (ac *Service) Update_Student_Details(update_student *models.StudentInfo, oldCourse string, oldName string, oldRollNumber string) error {

	course_old, err4 := ac.daos.GetCourseByName(oldCourse)
	course_new, err1 := ac.daos.GetCourseByName(update_student.ClassesEnrolled.CourseName)
	if err1 != nil {
		return fmt.Errorf("course not found - %s", update_student.ClassesEnrolled.CourseName)
	}
	if err4 != nil {
		return fmt.Errorf("course not found %s", oldCourse)
	}
	if _, err := strconv.ParseFloat(update_student.Name, 64); err == nil {
		return fmt.Errorf("name cant be number")
	}
	if update_student.Age <= 0 {
		return fmt.Errorf("invalid age")
	}
	existing_student, _ := ac.daos.GetStudentdetail(
		&models.StudentInfo{
			RollNumber: oldRollNumber,
			CourseId:   course_old.Id,
			Name:       oldName})

	if existing_student == nil {
		return fmt.Errorf("student details mismatched or does not exists")
	}

	if existing_student.Id == uuid.Nil {
		return fmt.Errorf("the student does not have %s registered", course_old.CourseName)
	}
	rcaNew, _ := ac.daos.GetStudentDetailsByRollNumberAndCourseId(update_student.RollNumber, course_new.Id)
	if rcaNew.Id != uuid.Nil && course_old.Id != course_new.Id {
		return fmt.Errorf("the student already has %s registered can't duplicate course ,please make course in url and body same", course_new.CourseName)
	}

	if course_old.Id != course_new.Id && rcaNew.Id == uuid.Nil {
		update_student.Id = existing_student.Id

	} else if rcaNew.Id != uuid.Nil {

		rcaOld, err := ac.daos.GetStudentDetailsByRollNumberAndCourseId(update_student.RollNumber, course_old.Id)

		if err != nil {
			return fmt.Errorf("student roll number not found %s", err.Error())
		}
		update_student.Id = rcaOld.Id
	} else {
		update_student.Id = existing_student.Id
	}

	sm, err2 := ac.daos.GetMarksByStudentId(update_student.Id)
	if err2 != nil {
		return err2
	}
	if update_student.ClassesEnrolled.Id == uuid.Nil {
		course_new, _ = ac.daos.GetCourseByName(update_student.ClassesEnrolled.CourseName)
		update_student.ClassesEnrolled.Id = course_new.Id
		update_student.CourseId = course_new.Id
	}
	update_student.MarksId = sm.Id
	sm.Marks = update_student.StudentMarks.Marks
	if sm.Marks > 100 {
		return fmt.Errorf("entered mark is beyond limit")
	}
	sm.Grade = ac.GenerateGradeForMarks(sm.Marks)
	sm.CourseName = update_student.ClassesEnrolled.CourseName
	sm.CourseId = update_student.CourseId
	update_student.StudentMarks = *sm

	err3 := ac.daos.UpdateStudentMarks(sm)
	if err3 != nil {
		return err3
	}
	err5 := ac.daos.UpdateClgStudent(update_student)
	if err5 != nil {
		return err5
	}
	if course_new.Id != course_old.Id {

		instructor_list_old, _ := ac.GetInstructorDetailWithSpecifics(models.InstructorDetails{CourseId: course_old.Id})
		for _, each_instructor := range instructor_list_old {
			err := ac.Update_Instructor_Info(each_instructor, models.InstructorDetails{Id: each_instructor.Id})
			if err != nil {
				return err
			}
		}

		instructor_list_new, _ := ac.GetInstructorDetailWithSpecifics(models.InstructorDetails{CourseId: course_new.Id})
		for _, each_instructor := range instructor_list_new {
			err1 := ac.Update_Instructor_Info(each_instructor, models.InstructorDetails{Id: each_instructor.Id})
			if err1 != nil {
				return err1
			}
		}
	}
	return nil

}

func (ac *Service) DeleteStudent(id uuid.UUID) error {

	student, err := ac.daos.GetStudentdetail(&models.StudentInfo{Id: id})
	if err != nil {
		return err
	}
	// for _, each_student := range student {
	err1 := ac.daos.DeleteStudentDaos(student.Id)
	if err1 != nil {
		return err1
	}
	// }
	// err1 := ac.daos.DeleteStudentDaos(student.Id)
	// if err1 != nil {
	// 	return err1
	// }
	return nil
}

func (ac *Service) UpdateStudentNameAge(existing_name, student_name string, age int64) error {
	si, err := ac.daos.GetStudentDetailsByName(existing_name)
	if err != nil {
		return err
	}
	for _, each_si := range *si {
		each_si.Name = student_name
		if age != 0 {
			each_si.Age = age
		}
		err := ac.daos.UpdateClgStudent(&each_si)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ac *Service) FetchStudentCourse(student_name string) (map[string]string, error) {
	course_list := make(map[string]string)
	course_list["student_name"] = student_name
	si, err := ac.daos.GetStudentDetailsByName(student_name)
	if err != nil {
		return nil, err
	}
	for index, each_si := range *si {
		each_si.ClassesEnrolled, err = ac.daos.GetCourseById(each_si.CourseId)
		if err != nil {
			return nil, err
		}

		studentmarks, err := ac.daos.GetMarksByStudentId(each_si.Id)
		each_si.StudentMarks = *studentmarks
		if err != nil {
			return nil, err
		}
		course_list[fmt.Sprintf("course_%d", index+1)] = fmt.Sprintf("%s -->Marks: %d Grade: %s", each_si.ClassesEnrolled.CourseName, each_si.StudentMarks.Marks, each_si.StudentMarks.Grade)
	}
	return course_list, nil
}

func (ac *Service) DeleteStudentCourseService(sn, cn string) (err error) {

	course_details, err1 := ac.daos.GetCourseByName(cn)
	if err1 != nil {
		return fmt.Errorf("course not found")
	}
	var student_detail models.StudentInfo
	student_detail.Name = sn
	student_detail.CourseId = course_details.Id

	_, err2 := ac.daos.GetStudentdetail(&student_detail)
	if err2 != nil {
		return err2
	}
	err = ac.daos.DeleteCourseForAStudent(sn, course_details.Id)
	if err != nil {
		return fmt.Errorf("failed to delete")
	}
	err5 := ac.daos.DeleteStudenetMarks(student_detail.MarksId)
	if err5 != nil {
		return fmt.Errorf("failed to delete")
	}
	return nil

}

func (s *Service) GetAllStudentSelectiveData() ([]*models.StudentSelectiveData, error) {

	ssd := []*models.StudentSelectiveData{}

	student_data, err := s.daos.RetrieveCollegeAdministration()
	if err != nil {
		return nil, err
	}
	for _, each_student_data := range student_data {
		ssd = append(ssd,
			&models.StudentSelectiveData{
				Name:       each_student_data.Name,
				RollNumber: each_student_data.RollNumber,
				Course:     each_student_data.ClassesEnrolled.CourseName,
			})

	}

	return ssd, nil

}
func (ac *Service) DeleteStudentSpecifics(st_req *models.StudentInfo) (err error) {

	course_details, err1 := ac.daos.GetCourseByName(st_req.ClassesEnrolled.CourseName)
	if err1 != nil {
		return fmt.Errorf("course not found")
	}

	st_details, err2 := ac.daos.GetStudentdetail(st_req)
	if err2 != nil {
		return err2
	}
	st_details.ClassesEnrolled = course_details

	err = ac.daos.DeleteStudentWithSpecifics(st_details)
	if err != nil {
		return fmt.Errorf("failed to delete")
	}
	return nil

}

func (ac *Service) UpdateStudentDetailsV2(update_student *models.StudentInfo) error {

	course_new, err1 := ac.daos.GetCourseByName(update_student.ClassesEnrolled.CourseName)
	if err1 != nil {
		return fmt.Errorf("course not found - %s", update_student.ClassesEnrolled.CourseName)
	}

	if _, err := strconv.ParseFloat(update_student.Name, 64); err == nil {
		return fmt.Errorf("name cant be number")
	}
	if update_student.Age <= 0 {
		return fmt.Errorf("invalid age")
	}
	existing_student, _ := ac.daos.GetStudentdetail(
		&models.StudentInfo{
			Id: update_student.Id})

	if existing_student == nil {
		return fmt.Errorf("student details mismatched or does not exists")
	}

	sm, err2 := ac.daos.GetMarksByStudentId(update_student.Id)
	if err2 != nil {
		return err2
	}
	if update_student.ClassesEnrolled.Id == uuid.Nil {
		course_new, _ = ac.daos.GetCourseByName(update_student.ClassesEnrolled.CourseName)
		update_student.ClassesEnrolled.Id = course_new.Id
		update_student.CourseId = course_new.Id
	}
	update_student.MarksId = sm.Id
	sm.Marks = update_student.StudentMarks.Marks
	if sm.Marks > 100 {
		return fmt.Errorf("entered mark is beyond limit")
	}
	sm.Grade = ac.GenerateGradeForMarks(sm.Marks)
	sm.CourseName = update_student.ClassesEnrolled.CourseName
	sm.CourseId = update_student.CourseId
	update_student.StudentMarks = *sm

	err3 := ac.daos.UpdateStudentMarks(sm)
	if err3 != nil {
		return err3
	}
	err5 := ac.daos.UpdateClgStudent(update_student)
	if err5 != nil {
		return err5
	}
	if course_new.Id != existing_student.CourseId {

		instructor_list_old, _ := ac.GetInstructorDetailWithSpecifics(models.InstructorDetails{CourseId: existing_student.CourseId})
		for _, each_instructor := range instructor_list_old {
			err := ac.Update_Instructor_Info(each_instructor, models.InstructorDetails{Id: each_instructor.Id})
			if err != nil {
				return err
			}
		}

		instructor_list_new, _ := ac.GetInstructorDetailWithSpecifics(models.InstructorDetails{CourseId: course_new.Id})
		for _, each_instructor := range instructor_list_new {
			err1 := ac.Update_Instructor_Info(each_instructor, models.InstructorDetails{Id: each_instructor.Id})
			if err1 != nil {
				return err1
			}
		}
	}
	return nil

}
