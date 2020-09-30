package action

import (
	"compass/internal/util"
	"compass/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

type Action struct {
	util.BaseModel
	WorkspaceId   string          `json:"workspaceId"`
	Nickname      string          `json:"nickname"`
	Type          string          `json:"type"`
	Description   string          `json:"description"`
	Configuration json.RawMessage `json:"configuration"`
	DeletedAt     time.Time       `json:"-"`
}

func (main Main) ParseAction(action io.ReadCloser) (Action, error) {
	var nAction *Action

	err := json.NewDecoder(action).Decode(&nAction)
	if err != nil {
		logger.Error(util.GeneralParseError, "ParseAction", err, action)
		return Action{}, err
	}

	return *nAction, nil
}

func (main Main) ValidateAction(action Action) []util.ErrorUtil {
	ers := make([]util.ErrorUtil, 0)
	needConfigValidation := true

	if action.Nickname == "" {
		ers = append(ers, util.ErrorUtil{Field: "nickname", Error: errors.New("action nickname is required").Error()})
	}

	if action.Configuration == nil || len(action.Configuration) == 0 {
		ers = append(ers, util.ErrorUtil{Field: "configuration", Error: errors.New("action configuration is required").Error()})
		needConfigValidation = false
	}

	if action.WorkspaceId == "" {
		ers = append(ers, util.ErrorUtil{Field: "workspaceId", Error: errors.New("workspaceId is required").Error()})
	}

	if action.Type == "" {
		ers = append(ers, util.ErrorUtil{Field: "type", Error: errors.New("action type is required").Error()})
		needConfigValidation = false
	}

	if needConfigValidation {
		ers = append(ers, main.validateActionConfig(action.Type, action.Configuration)...)
	}

	return ers
}

func (main Main) validateActionConfig(actionType string, actionConfiguration json.RawMessage) []util.ErrorUtil {
	ers := make([]util.ErrorUtil, 0)

	plugin, err := main.pluginRepo.GetPluginBySrc(actionType)
	if err != nil {
		logger.Error("error finding plugin", "ValidateActionConfig", err, actionType)
		return []util.ErrorUtil{{Field: "type", Error: errors.New("action type is invalid").Error()}}
	}

	pluginErrs, err := plugin.Lookup("ValidateActionConfiguration")
	if err != nil {
		logger.Error("error looking up for plugin", "ValidateActionConfig", err, fmt.Sprintf("%+v", plugin))
		return []util.ErrorUtil{{Field: "type", Error: errors.New("action type is invalid").Error()}}
	}

	configErs := pluginErrs.(func(actionConfig []byte) []error)(actionConfiguration)
	if len(configErs) > 0 {
		for _, err := range configErs {
			ers = append(ers, util.ErrorUtil{Field: "configuration", Error: err.Error()})
		}
	}

	return ers
}

func (main Main) FindActionById(id string) (Action, error) {
	action := Action{}
	db := main.db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&action)
	if db.Error != nil {
		logger.Error(util.FindActionError, "FindActionById", db.Error, "Id = "+id)
		return Action{}, db.Error
	}

	return action, nil
}

func (main Main) FindAllActionsByWorkspace(workspaceID string) ([]Action, error) {
	var actions []Action

	db := main.db.Set("gorm:auto_preload", true).Where("workspace_id = ? and deleted_at is null", workspaceID).Find(&actions)
	if db.Error != nil {
		logger.Error(util.FindActionError, "FindAllActionsByWorkspace", db.Error, actions)
		return []Action{}, db.Error
	}
	return actions, nil
}

func (main Main) SaveAction(action Action) (Action, error) {
	db := main.db.Create(&action)
	if db.Error != nil {
		logger.Error(util.SaveActionError, "SaveAction", db.Error, action)
		return Action{}, db.Error
	}
	return action, nil
}

func (main Main) DeleteAction(id string) error {
	db := main.db.Model(&Action{}).Where("id = ?", id).Delete(&Action{})
	if db.Error != nil {
		logger.Error(util.DeleteActionError, "DeleteAction", db.Error, "Id = "+id)
		return db.Error
	}
	return nil
}
