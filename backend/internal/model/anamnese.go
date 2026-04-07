package model

import (
	"time"

	"github.com/google/uuid"
)

//identificação do paciente
type IdentificationSection struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	BirthDate   string `json:"birth_date"`
	Sex         string `json:"sex"`
	Naturality  string `json:"naturality"`
	School      string `json:"school"`
	Grade       string `json:"grade"`
	Classroom   string `json:"classroom"`
}

//identificação da familia
type FamilyDataSection struct {
	FatherName        *string `json:"father_name"`
	FatherAge         *int    `json:"father_age"`
	FatherProfession  *string `json:"father_profession"`
	MotherName        string `json:"mother_name"`
	MotherAge         int    `json:"mother_age"`
	MotherProfession  *string `json:"mother_profession"`
	Guardians         string `json:"guardians"`
}

//informações sobre a gravidez
type PregnancyBirthSection struct {
	PlannedPregnancy bool   `json:"planned_pregnancy"`
	Complications    string `json:"complications"`
	BirthType        string `json:"birth_type"`
	Premature        bool   `json:"premature"`
	Apgar            string `json:"apgar"`
}

//amamentação
type BreastfeedingSection struct {
	Breastfed      bool   `json:"breastfed"`
	Duration       string `json:"duration"`
	Difficulties   string `json:"difficulties"`
}

//seção de alimentação
type FeedingSection struct {
	EatsWell         bool   `json:"eats_well"`
	FoodSelectivity  string `json:"food_selectivity"`
	Difficulties     string `json:"difficulties"`
}

//seção de indepedência do paciente
type IndependenceSection struct {
	BathesAlone     bool `json:"bathes_alone"`
	DressesAlone    bool `json:"dresses_alone"`
	HygieneAlone    bool `json:"hygiene_alone"`
}

//seção da saúde do paciente
type HealthSection struct {
	Diseases     string `json:"diseases"`
	Medications  string `json:"medications"`
	Therapies    string `json:"therapies"`
}

//seção de sono do paciente
type SleepSection struct {
	SleepsWell    bool   `json:"sleeps_well"`
	WakesUpNight  bool   `json:"wakes_up_night"`
	SleepRoutine  string `json:"sleep_routine"`
}

//seção de coordenação do paciente
type PsychomotorSection struct {
	Coordination string `json:"coordination"`
	Balance      string `json:"balance"`
	Laterality   string `json:"laterality"`
}

//seção de linguagem do paciente
type LanguageSection struct {
	SpeaksWell       bool   `json:"speaks_well"`
	Communication    string `json:"communication"`
	Difficulties     string `json:"difficulties"`
}

//seção de socialização do paciente
type SocializationSection struct {
	HasFriends     bool   `json:"has_friends"`
	Interaction    string `json:"interaction"`
	Behavior       string `json:"behavior"`
}


//desenvolvimento familiar
type FamilyEnvironmentSection struct {
	HomeRoutine string `json:"home_routine"`
	Conflicts   string `json:"conflicts"`
	Relationship string `json:"relationship"`
}

//histórico(genético) familiar
type FamilyHistorySection struct {
	LearningIssues string `json:"learning_issues"`
	Diseases       string `json:"diseases"`
}

//atividades extras
type ExtracurricularSection struct {
	Activities *string `json:"activities"`
	Frequency  *string `json:"frequency"`
}

//seção escolar do paciente
type SchoolingSection struct {
	Performance   string `json:"performance"`
	Difficulties  string `json:"difficulties"`
	TeacherNotes  string `json:"teacher_notes"`
}

//personalidade do paciente
type PersonalitySection struct {
	Traits string `json:"traits"`
	Mood   string `json:"mood"`
}

//gostos do paciente
type PreferencesSection struct {
	Likes    string `json:"likes"`
	Dislikes string `json:"dislikes"`
}



type Anamnesis struct {
	ID        uuid.UUID `json:"id"`
	PatientID uuid.UUID `json:"patient_id"`

	Identification    IdentificationSection    `json:"identification"`
	FamilyData        FamilyDataSection        `json:"family_data"`
	PregnancyBirth    PregnancyBirthSection    `json:"pregnancy_birth"`
	Breastfeeding     BreastfeedingSection     `json:"breastfeeding"`
	Feeding           FeedingSection           `json:"feeding"`
	Independence      IndependenceSection      `json:"independence"`
	Health            HealthSection            `json:"health"`
	Sleep             SleepSection             `json:"sleep"`
	Psychomotor       PsychomotorSection       `json:"psychomotor"`
	Language          LanguageSection          `json:"language"`
	Socialization     SocializationSection     `json:"socialization"`
	FamilyEnvironment FamilyEnvironmentSection `json:"family_environment"`
	FamilyHistory     FamilyHistorySection     `json:"family_history"`
	Extracurricular   ExtracurricularSection   `json:"extracurricular"`
	Schooling         SchoolingSection         `json:"schooling"`
	Personality       PersonalitySection       `json:"personality"`
	Preferences       PreferencesSection       `json:"preferences"`

	Routine       	string `json:"routine"`
	AdditionalInfo 	string `json:"additional_info"`

	Version         int        `json:"version"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}