package coverageoutput

import (
	"fmt"
	"log"

	"tests-coverage-tool/tool/config"
	"tests-coverage-tool/tool/coverageinupt"
	"tests-coverage-tool/tool/reflection"
	"tests-coverage-tool/tool/utils"

	"github.com/jhump/protoreflect/desc"
	"github.com/samber/lo"
)

type OutputCoverageClient struct {
	reflectionClient    reflection.GRPCReflectionClient
	inputCoverageClient coverageinupt.InputCoverageClient
}

func NewOutputCoverageClient(
	reflectionClient *reflection.GRPCReflectionClient,
	inputCoverageClient *coverageinupt.InputCoverageClient,
) (*OutputCoverageClient, error) {
	if reflectionClient == nil {
		return nil, fmt.Errorf("reflection client is nil")
	}

	if inputCoverageClient == nil {
		return nil, fmt.Errorf("input coverage client is nil")
	}

	return &OutputCoverageClient{
		reflectionClient:    *reflectionClient,
		inputCoverageClient: *inputCoverageClient,
	}, nil
}

func (c *OutputCoverageClient) getParameterCoverage(parameter string, resultParameters []string) ParameterCoverage {
	return ParameterCoverage{
		Covered:    lo.Contains(lo.Uniq(resultParameters), parameter),
		Parameter:  parameter,
		TotalCases: lo.Count(resultParameters, parameter),
	}
}

func (c *OutputCoverageClient) getMethodCoverage(descriptor *desc.MethodDescriptor, logicalService string) MethodCoverage {
	method := descriptor.GetName()
	requestParameters := reflection.GetFieldNamesFromMessageDescriptor(descriptor.GetInputType())
	responseParameters := reflection.GetFieldNamesFromMessageDescriptor(descriptor.GetOutputType())

	filters := coverageinupt.ResultsFilters{Method: method, LogicalService: logicalService}
	resultsRequestParameters := c.inputCoverageClient.GetRequestParameters(filters)
	resultsResponseParameters := c.inputCoverageClient.GetResponseParameters(filters)

	requestParametersCoverage := lo.Map(requestParameters, func(item string, _ int) ParameterCoverage {
		return c.getParameterCoverage(item, resultsRequestParameters)
	})
	responseParametersCoverage := lo.Map(responseParameters, func(item string, _ int) ParameterCoverage {
		return c.getParameterCoverage(item, resultsResponseParameters)
	})

	return MethodCoverage{
		Method:     method,
		Covered:    c.inputCoverageClient.IsMethodCoveted(filters),
		TotalCases: len(c.inputCoverageClient.FilterResults(filters)),
		RequestCoverage: MethodRequestCoverage{
			Name:                   descriptor.GetInputType().GetFullyQualifiedName(),
			TotalParameters:        len(requestParameters),
			ParametersCoverage:     requestParametersCoverage,
			TotalCoveredParameters: len(resultsRequestParameters),
		},
		ResponseCoverage: MethodRequestCoverage{
			Name:                   descriptor.GetOutputType().GetFullyQualifiedName(),
			TotalParameters:        len(responseParameters),
			ParametersCoverage:     responseParametersCoverage,
			TotalCoveredParameters: len(resultsResponseParameters),
		},
	}
}

func (c *OutputCoverageClient) getLogicalServiceCoverage(logicalService string) (LogicalServiceCoverage, error) {
	filters := coverageinupt.ResultsFilters{LogicalService: logicalService}
	resultsMethods := c.inputCoverageClient.GetUniqueMethods(filters)

	reflectionMethodsDescriptors, err := c.reflectionClient.GetServiceMethodsDescriptors(logicalService)
	if err != nil {
		return LogicalServiceCoverage{}, err
	}
	reflectionMethods := lo.Map(reflectionMethodsDescriptors, func(item *desc.MethodDescriptor, _ int) string {
		return item.GetName()
	})

	return LogicalServiceCoverage{
		Methods: lo.Map(reflectionMethodsDescriptors, func(item *desc.MethodDescriptor, _ int) MethodCoverage {
			return c.getMethodCoverage(item, logicalService)
		}),
		TotalCoverage:  getCoveragePercent(reflectionMethods, resultsMethods),
		LogicalService: logicalService,
	}, nil
}

func (c *OutputCoverageClient) getLogicalServiceCoverages() ([]LogicalServiceCoverage, error) {
	logicalServices, err := c.reflectionClient.GetServices()
	if err != nil {
		return nil, err
	}

	coverages := make([]LogicalServiceCoverage, len(logicalServices))
	for index, logicalService := range logicalServices {
		if coverage, err := c.getLogicalServiceCoverage(logicalService); err == nil {
			coverages[index] = coverage
		}
	}

	return coverages, nil
}

func (c *OutputCoverageClient) SaveResults(resultsDir string) ([]LogicalServiceCoverage, error) {
	log.Printf("Starting to save logical service coverages into results folder")

	logicalServicesCoverage, err := c.getLogicalServiceCoverages()
	if err != nil {
		return nil, err
	}

	if resultsDir == "" {
		log.Println("Env variable 'TESTS_COVERAGE_INPUT_RESULTS_DIR' empty, skipping")
		return logicalServicesCoverage, nil
	}

	if err = utils.SaveJSONFile(logicalServicesCoverage, resultsDir, config.StateLogicalServicesCoverageJSON); err != nil {
		return nil, err
	}

	return logicalServicesCoverage, nil
}
