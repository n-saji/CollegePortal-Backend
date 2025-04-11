package service

import (
	"CollegeAdministration/config"
	"CollegeAdministration/models"
	"CollegeAdministration/utils"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (ac *Service) InsertInstructor(account_id string, iid *models.InstructorDetails) (uuid.UUID, error) {
	if _, err := strconv.Atoi(iid.InstructorName); err == nil {
		return uuid.Nil, errors.New("name can't be number")
	}
	if _, err := strconv.Atoi(iid.Department); err == nil {
		return uuid.Nil, errors.New("departmment can't be number")
	}
	cn, err1 := ac.daos.GetCourseByName(iid.CourseName)
	if err1 != nil {
		return uuid.Nil, fmt.Errorf("course not available")
	}
	iid.CourseId = cn.Id
	cd_exist, _ := ac.daos.GetInstructor(iid)
	if cd_exist.Id != uuid.Nil {
		return uuid.Nil, fmt.Errorf("instructor exits")
	}
	ok, err2 := ac.ValidateInstructorDetails(iid)

	if !ok {
		return uuid.Nil, err2
	}
	iid.Id = uuid.New()
	StudentList, _ := ac.daos.RetrieveCollegeAdministration()
	var student_with_course []models.StudentInfo
	for _, each_student := range StudentList {
		if each_student.CourseId == iid.CourseId {
			student_with_course = append(student_with_course, *each_student)
		}
	}
	iid.Info.StudentsList = student_with_course
	err := ac.daos.InsertInstructorDetails(iid)
	if err != nil {
		log.Println("Error while inserting details")
		return uuid.Nil, err
	}
	account := &models.Account{}
	account.Id = iid.Id
	account.Info.Credentials.Id = iid.Id
	account.Name = iid.InstructorName
	account.Type = config.AccountTypeInstructor
	account.Verified = false

	err3 := ac.daos.CreateAccount(account)
	if err3 != nil {
		log.Println("error storing in account")
		return uuid.Nil, err3
	}

	go utils.StoreMessages("New Instructor", account.Name, config.AccountTypeInstructor, account_id)
	return iid.Id, nil
}

func (ac *Service) GetInstructorDetails() ([]*models.InstructorDetails, error) {

	id, err := ac.daos.GetAllInstructor()
	for _, eachId := range id {
		course, err := ac.daos.GetCourseById(eachId.CourseId)
		if err != nil {
			return nil, err
		}
		eachId.CourseName = course.CourseName
		eachId.ClassesEnrolled, _ = ac.daos.GetCourseByName(eachId.CourseName)
	}
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *Service) GetInstructorDetailsWithConditions(order_clause string) ([]*models.InstructorDetails, error) {

	id, err := s.daos.GetAllInstructorOrderByCondition(order_clause)
	if err != nil {
		return nil, err
	}

	for i := range id {
		id[i].ClassesEnrolled, err = s.daos.GetCourseById(id[i].CourseId)
		if err != nil {
			return nil, err
		}
		id[i].CourseName = id[i].ClassesEnrolled.CourseName
	}

	return id, nil
}

func (ac *Service) StoreInstructoLogindetails(id uuid.UUID, emailid, password string) error {

	var credentials models.InstructorLogin

	crypted_password, err2 := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err2 != nil {
		return fmt.Errorf("error parsing password")
	}
	credentials.Id = id
	credentials.EmailId = emailid
	credentials.Password = string(crypted_password)

	if id == uuid.Nil {
		return fmt.Errorf("uuid cant be null")
	}
	err := ac.daos.CheckIDPresent(credentials.Id)
	if err != nil {
		return err
	}
	err = ac.daos.CreateInstructorLogin(credentials)
	if err != nil {
		return err
	}
	iid, err := ac.daos.GetInstructor(&models.InstructorDetails{Id: id})
	if err != nil {
		return err
	}
	account := &models.Account{}
	account.Id = iid.Id
	account.Info.Credentials = models.InstructorLogin{Id: iid.Id,
		EmailId:  emailid,
		Password: string(crypted_password)}
	account.Name = iid.InstructorName
	account.Type = config.AccountTypeInstructor
	account.Verified = false

	err3 := ac.daos.CreateAccount(account)
	if err3 != nil {
		log.Println("error storing in account instructor")
		return err3
	}

	go ac.GenerateOTPAndStore(emailid)

	return nil
}

func (s *Service) DeleteInstructor(name string) error {

	id, err := s.daos.GetInstructorWithName(name)
	if id.Id == uuid.Nil {
		return fmt.Errorf("instructor not found")
	}
	if err != nil {
		return err
	}

	err3 := s.daos.DeleteMessageByAccountId(id.Id)
	if err3 != nil {
		return err3
	}

	err1 := s.daos.DeleteInstructorLogin(id.Id)
	if err1 != nil {
		return err1
	}

	err2 := s.daos.DeleteInstructor(name)
	if err2 != nil {
		return err2
	}
	return nil
}

func (s *Service) Update_Instructor_Info(req_id *models.InstructorDetails, cond *models.InstructorDetails) error {
	list_details, err := s.GetInstructorDetailWithSpecifics(cond)
	if err != nil {
		return fmt.Errorf(fmt.Sprint("fetching error " + err.Error()))
	}
	if len(list_details) == 0 {
		return fmt.Errorf("no instructor found with given details")
	}
	if req_id.CourseName != "" {
		status := s.daos.CheckCourse(req_id.CourseName)
		if !status {
			return fmt.Errorf("course does not exits")
		}
		course_details, _ := s.daos.GetCourseByName(req_id.CourseName)
		req_id.CourseId = course_details.Id
	}

	ok, err2 := s.ValidateInstructorDetails(req_id)

	if !ok {
		return err2
	}

	StudentList, _ := s.daos.RetrieveCollegeAdministration()
	var student_with_course []models.StudentInfo
	for _, each_student := range StudentList {
		if each_student.CourseId == req_id.CourseId {
			student_with_course = append(student_with_course, *each_student)
		}
	}
	req_id.Info.StudentsList = student_with_course
	err1 := s.daos.UpdateInstructorInfo(req_id, cond)

	if err1 != nil {
		return err1
	}

	return nil
}

func (s *Service) Update_Instructor(req_id *models.InstructorDetails, cond *models.InstructorDetails) error {

	list_details, err := s.GetInstructorDetailWithSpecifics(cond)
	if err != nil {
		return fmt.Errorf(fmt.Sprint("fetching error " + err.Error()))
	}
	if len(list_details) == 0 {
		return fmt.Errorf("no instructor found with given details")
	}
	if req_id.CourseName != "" {
		status := s.daos.CheckCourse(req_id.CourseName)
		if !status {
			return fmt.Errorf("course does not exits")
		}
		course_details, _ := s.daos.GetCourseByName(req_id.CourseName)
		req_id.CourseId = course_details.Id
	}

	ok, err2 := s.ValidateInstructorDetails(req_id)

	if !ok {
		return err2
	}

	StudentList, _ := s.daos.RetrieveCollegeAdministration()
	var student_with_course []models.StudentInfo
	for _, each_student := range StudentList {
		if each_student.CourseId == req_id.CourseId {
			student_with_course = append(student_with_course, *each_student)
		}
	}
	req_id.Info.StudentsList = student_with_course
	err1 := s.daos.UpdateInstructor(req_id, cond)

	if err1 != nil {
		return err1
	}

	return nil

}
func (s *Service) GetInstructorDetailWithSpecifics(req *models.InstructorDetails) ([]*models.InstructorDetails, error) {

	id_list, err := s.daos.RetieveInstructorDetailsWithCondition(req)
	if err != nil {
		return nil, fmt.Errorf("error %s", err.Error())
	}

	return id_list, nil

}

func (s *Service) DeleteInstructorWithConditions(aid string, id_condition *models.InstructorDetails) error {

	if aid == "" {
		return fmt.Errorf("account id not found")
	}
	if id_condition != nil {
		if id_condition.Id != uuid.Nil && id_condition.Id.String() == aid {
			return fmt.Errorf("account id and instructor id are same")
		}
	}

	id_list, err := s.daos.RetieveInstructorDetailsWithCondition(id_condition)
	if err != nil {
		return err
	}

	for _, each_id := range id_list {
		if each_id.Id == uuid.Nil {
			return fmt.Errorf("instructor not found")
		}

		err = s.daos.DeleteMessageByAccountId(each_id.Id)
		if err != nil {
			return fmt.Errorf("error deleting message %s", err.Error())
		}

		err = s.daos.DeleteInstructorLogin(each_id.Id)
		if err != nil {
			return fmt.Errorf("error deleting instructor login %s", err.Error())
		}

		err = s.daos.DeleteTokenByAccountId(each_id.Id)
		if err != nil {
			return fmt.Errorf("error deleting token %s", err.Error())
		}

		err = s.daos.DeleteInstructorWithConditions(each_id)
		if err != nil {
			return fmt.Errorf("error deleting instructor %s", err.Error())
		}
	}

	return nil
}

func (s *Service) GetInstructorIDWithEmail(email string) (string, error) {

	instructor_id, err := s.daos.GetIDUsingEmail(email)
	if err != nil {
		return "", fmt.Errorf("failed fetching :" + err.Error())
	}
	return instructor_id, nil
}

func (s *Service) GetInstructorNamewithId(id string) (*models.InstructorDetails, error) {
	iid := &models.InstructorDetails{}
	id_uuid, err1 := uuid.Parse(id)
	if err1 != nil {
		return nil, err1
	}
	iid.Id = id_uuid
	i_details, err := s.daos.GetInstructor(iid)
	if err != nil {
		return nil, err
	}
	return i_details, nil
}

func (s *Service) ViewinstructorProfile(i_id string) (*models.InstructorProfile, error) {
	Profile := &models.InstructorProfile{}
	i_id_parsed, err := uuid.Parse(i_id)
	if err != nil {
		return nil, err
	}
	instructor_detail, err := s.daos.GetInstructor(&models.InstructorDetails{Id: i_id_parsed})
	if err != nil {
		return nil, err
	}
	credentials, err1 := s.daos.FetchCredentialsUsingID(i_id_parsed)
	if err1 != nil {
		return nil, err1
	}
	Profile.Code = instructor_detail.InstructorCode
	Profile.Name = instructor_detail.InstructorName
	Profile.Department = instructor_detail.Department
	Profile.CourseList = instructor_detail.CourseName
	Profile.Credentials = models.InstructorLogin{Id: credentials.Id,
		EmailId:  credentials.EmailId,
		Password: credentials.Password}

	return Profile, nil
}

func (s *Service) UpdateInstructorCredentials(cred *models.InstructorLogin) error {

	var crypted_password []byte
	var err1 error

	if cred.EmailId == "" {
		existing_credentials, err := s.daos.FetchCredentialsUsingID(cred.Id)
		if err != nil {
			return err
		}
		cred.EmailId = existing_credentials.EmailId
	} else {
		err2 := s.ValidateEmail(cred.EmailId)
		if err2 != nil {
			return err2
		}
		exists, err := s.CheckEmailExist(cred.EmailId)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("email already exists")
		}

	}

	if cred.Password != "" {
		err2 := s.ValidatePassword(cred.Password)
		if err2 != nil {
			return err2
		}
		crypted_password, err1 = bcrypt.GenerateFromPassword([]byte(cred.Password), 10)
		if err1 != nil {
			return err1
		}
		cred.Password = string(crypted_password)
	} else {
		existing_credentials, err := s.daos.FetchCredentialsUsingID(cred.Id)
		if err != nil {
			return err
		}
		cred.Password = existing_credentials.Password
	}

	err := s.daos.UpdateCredentials(cred)
	if err != nil {
		return err
	}
	return nil
}
