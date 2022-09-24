package xplorersbot

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
)

type SsmApi struct{}

func (ssmApi SsmApi) GetParameter(ctx context.Context,
	params *ssm.GetParameterInput,
	optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {

	parameter := &types.Parameter{Value: aws.String("some-parameter-value")}

	output := &ssm.GetParameterOutput{
		Parameter: parameter,
	}

	return output, nil
}

func TestGetParameter(t *testing.T) {
	assert := require.New(t)
	api := &SsmApi{}
	parameterPath := "/some/parameter/path"

	input := &ssm.GetParameterInput{
		Name: &parameterPath,
	}

	result, err := GetParameter(context.Background(), *api, input)
	if err != nil {
		t.Log("Unable to fetch parameter")
		t.Log(err)
		return
	}

	t.Log("Parameter value: " + *result.Parameter.Value)
	assert.Equal(*result.Parameter.Value, "some-parameter-value")
}
