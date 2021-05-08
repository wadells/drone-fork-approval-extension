package plugin

import (
	"context"

	"github.com/drone/drone-go/plugin/validator"
	"github.com/sirupsen/logrus"
)

// New returns a new validator plugin.
func New(log logrus.FieldLogger) validator.Plugin {
	return &plugin{log: log}
}

type plugin struct {
	log logrus.FieldLogger
}

func (p *plugin) Validate(ctx context.Context, req *validator.Request) error {
	switch req.Build.Event {
	case "push", "tag": // triggered by folks with write access to the repo, therefore trusted
		p.log.Debugf("%s build ignored", req.Build.Event)
		return nil
	case "cron", "promote", "rollback", "custom": // triggered by folks with write access in drone
		p.log.Debugf("%s build ignored", req.Build.Event)
		return nil
	case "pull_request": // may be triggered by folks without write access, needs approval
		break
	default: // unknown new event, fail secure
		p.log.Warnf("%s build unrecognized", req.Build.Event)
		return validator.ErrBlock
	}

	sourceRepo := req.Build.Fork
	targetRepo := req.Repo.Slug
	if sourceRepo != targetRepo { // then the PR is coming from a fork
		p.log.WithFields(logrus.Fields{"source": sourceRepo, "target": targetRepo}).Infof("%s needs approval", req.Build.Link)
		return validator.ErrBlock
	} else {
		p.log.WithFields(logrus.Fields{"source": sourceRepo, "target": targetRepo}).Infof("%s approved", req.Build.Link)
		return nil
	}
}
