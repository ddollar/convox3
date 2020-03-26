package resolver

import (
	"context"
	"fmt"

	"github.com/convox/console/api/model"
)

func cmodel(ctx context.Context) (model.Interface, error) {
	m, ok := ctx.Value("model").(model.Interface)
	if !ok {
		return nil, fmt.Errorf("invalid model")
	}

	return m, nil
}

func cuid(ctx context.Context) (string, error) {
	uid, ok := ctx.Value("uid").(string)
	if !ok {
		return "", fmt.Errorf("invalid uid")
	}

	return uid, nil
}
