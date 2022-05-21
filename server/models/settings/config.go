package settings

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config struct for webapp config
type Config struct {
	Messages struct {
		// MsgIncorrectInput appears when user's input was parsed unsuccessfully
		MsgIncorrectInput string `yaml:"msg_incorrect_input"`

		// MsgQuestionRepeat appears when user asks to repeat last question
		MsgQuestionRepeat string `yaml:"msg_question_repeat"`

		// MsgStartQuizOver appears when user asks to restart quiz
		MsgStartQuizOver string `yaml:"msg_start_quiz_over"`

		// MsgFinishQuiz appears after user completes a quiz
		MsgFinishQuiz string `yaml:"msg_finish_quiz"`

		// MsgHelp appears when user ask for help
		MsgHelp string `yaml:"msg_help"`

		// MsgAlreadyLoggedIn appears when user tries to log in while being logged in already
		MsgAlreadyLoggedIn string `yaml:"msg_already_logged_in"`

		// MsgRegistrationPrompt appears when user logs in and their session_id is missing in database
		MsgRegistrationPrompt string `yaml:"msg_registration_prompt"`

		// MsgWelcomeAfterRegistration appears after user was registrated and logged in
		MsgWelcomeAfterRegistration string `yaml:"msg_welcome_after_registration"`

		// MsgWelcomeAfterLogin appears after already registered user logs in
		MsgWelcomeAfterLogin string `yaml:"msg_welcome_after_login"`

		// MsgGoodbye appears after user quits skill
		MsgGoodbye string `yaml:"msg_goodbye"`
	} `yaml:"messages"`

	Secrets struct {
		DBHost     string `yaml:"db_host"`
		DBPort     string `yaml:"db_port"`
		DBName     string `yaml:"db_name"`
		DBUser     string `yaml:"db_user"`
		DBPassword string `yaml:"db_pass"`
		DBSSL      string `yaml:"db_ssl"` // disable, prefer etc.
	}
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
