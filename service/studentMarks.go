package service

import (
	"CollegeAdministration/models"
	"fmt"

	"github.com/google/uuid"
)

func (ac *Service) InsertStudentIdInToMarksTable(cv *models.StudentInfo) (*models.StudentMarks, error) {

	var sm models.StudentMarks

	sm.Id = uuid.New()
	sm.CourseId = cv.ClassesEnrolled.Id
	sm.CourseName = cv.ClassesEnrolled.CourseName
	sm.StudentId = cv.Id
	sm.Grade = "nil"
	err := ac.daos.CreateStudentMarks(&sm)
	if err != nil {
		return nil, err
	}
	return &sm, nil
}

func (s *Service) GetAllStudentsMarksForGivenCourse(course_name string) (*models.StudentsMarksForCourse, error) {

	smfc := &models.StudentsMarksForCourse{}
	smfc.StudentNameMark = map[string]int64{}
	smfc.Ranking = map[int64]string{}
	var err error
	course_model, err1 := s.daos.GetCourseByName(course_name)
	if err1 != nil {
		return nil, fmt.Errorf("no course exits %s", err1.Error())
	}
	smfc.StudentId, err = s.daos.GetAllStudentsIDForACourse(course_model.Id)
	if err != nil {
		return nil, err
	}
	if len(smfc.StudentId) == 0 {
		return nil, fmt.Errorf("no student data exists")
	}
	smfc.Course_name = course_model.CourseName

	for index, each_student := range smfc.StudentId {
		student_id, err := uuid.Parse(each_student)
		if err != nil {
			return nil, err
		}
		student_model, err1 := s.daos.GetStudentdetail(&models.StudentInfo{Id: student_id, CourseId: course_model.Id})
		if err1 != nil {
			return nil, err1
		}
		marks_model, err2 := s.daos.GetMarksByStudentId(student_id)
		if err2 != nil {
			return nil, err2
		}
		student_model.StudentMarks = *marks_model
		smfc.StudentNameMark[student_model.Name] = student_model.StudentMarks.Marks
		smfc.Ranking[int64(index+1)] = student_model.Name

	}

	return smfc, nil

}

func (s *Service) GenerateGradeForMarks(marks int64) string {

	if marks > 90 {
		return "S"
	} else if marks <= 90 && marks > 80 {
		return "A+"
	} else if marks <= 80 && marks > 70 {
		return "A"
	} else if marks <= 70 && marks > 60 {
		return "B+"
	} else if marks <= 60 && marks > 50 {
		return "B"
	} else if marks <= 50 && marks > 40 {
		return "C"
	} else {
		return "Fail"
	}

}
