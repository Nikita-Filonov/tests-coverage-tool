package coverageoutput

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/samber/lo"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverage"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverageinupt"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/reflection"
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

func (c *OutputCoverageClient) getMethodCoverage(descriptor *desc.MethodDescriptor) MethodCoverage {
	var (
		option                     = descriptor.GetMethodOptions()
		filters                    = coverageinupt.ResultsFilters{FilterByFullMethod: descriptor.GetFullyQualifiedName()}
		results                    = c.inputCoverageClient.FilterResults(filters)
		covered                    = len(results) > 0
		totalCases                 = len(results)
		actualRequestParameters    = c.inputCoverageClient.GetMergedRequestParameters(filters)
		actualResponseParameters   = c.inputCoverageClient.GetMergedResponseParameters(filters)
		expectedRequestParameters  = reflection.BuildExpectedResultParameters(descriptor.GetInputType().UnwrapMessage())
		expectedResponseParameters = reflection.BuildExpectedResultParameters(descriptor.GetOutputType().UnwrapMessage())
	)

	var (
		requestTotalParameters         = coverage.GetTotalResultParameters(expectedRequestParameters)
		responseTotalParameters        = coverage.GetTotalResultParameters(expectedResponseParameters)
		requestTotalCoveredParameters  = coverage.GetTotalCoveredResultParameters(actualRequestParameters)
		responseTotalCoveredParameters = coverage.GetTotalCoveredResultParameters(actualResponseParameters)
		requestParametersCoverage      []coverage.ResultParameters
		responseParametersCoverage     []coverage.ResultParameters
	)
	if covered {
		requestParametersCoverage = coverage.MergeResultParameters(expectedRequestParameters, actualRequestParameters)
		responseParametersCoverage = coverage.MergeResultParameters(expectedResponseParameters, actualResponseParameters)
	}

	return MethodCoverage{
		Method:     descriptor.GetName(),
		Covered:    covered,
		TotalCases: totalCases,
		Deprecated: option.GetDeprecated(),
		RequestCoverage: MethodRequestCoverage{
			Name:                   descriptor.GetInputType().GetFullyQualifiedName(),
			TotalCoverage:          getRequestCoveragePercent(covered, requestTotalParameters, requestTotalCoveredParameters),
			TotalParameters:        requestTotalParameters,
			ParametersCoverage:     requestParametersCoverage,
			TotalCoveredParameters: requestTotalCoveredParameters,
		},
		ResponseCoverage: MethodRequestCoverage{
			Name:                   descriptor.GetOutputType().GetFullyQualifiedName(),
			TotalCoverage:          getRequestCoveragePercent(covered, responseTotalParameters, responseTotalCoveredParameters),
			TotalParameters:        responseTotalParameters,
			ParametersCoverage:     responseParametersCoverage,
			TotalCoveredParameters: responseTotalCoveredParameters,
		},
	}
}

func (c *OutputCoverageClient) getLogicalServiceCoverage(logicalService string) (LogicalServiceCoverage, error) {
	filters := coverageinupt.ResultsFilters{FilterByLogicalService: logicalService}
	resultsMethods := c.inputCoverageClient.GetUniqueMethods(filters)

	reflectionMethodsDescriptors, err := c.reflectionClient.GetServiceMethodsDescriptors(logicalService)
	if err != nil {
		return LogicalServiceCoverage{}, err
	}

	reflectionMethods, err := c.reflectionClient.GetServiceMethods(logicalService)
	if err != nil {
		return LogicalServiceCoverage{}, err
	}

	return LogicalServiceCoverage{
		Methods:        lo.Map(reflectionMethodsDescriptors, func(item *desc.MethodDescriptor, _ int) MethodCoverage { return c.getMethodCoverage(item) }),
		TotalCoverage:  getCoveragePercent(reflectionMethods, resultsMethods),
		LogicalService: logicalService,
	}, nil
}

func (c *OutputCoverageClient) GetLogicalServiceCoverages() ([]LogicalServiceCoverage, error) {
	logicalServices, err := c.reflectionClient.GetServices()
	if err != nil {
		return nil, err
	}

	coverages := make([]LogicalServiceCoverage, len(logicalServices))
	for index, logicalService := range logicalServices {
		if serviceCoverage, err := c.getLogicalServiceCoverage(logicalService); err == nil {
			coverages[index] = serviceCoverage
		}
	}

	return coverages, nil
}

func (c *OutputCoverageClient) GetServiceCoverage() (ServiceCoverage, error) {
	logicalServices, err := c.reflectionClient.GetServices()
	if err != nil {
		return ServiceCoverage{}, err
	}

	var (
		resultsMethods    []string
		reflectionMethods []string
	)
	for _, logicalService := range logicalServices {
		filters := coverageinupt.ResultsFilters{FilterByLogicalService: logicalService}
		resultsMethods = append(resultsMethods, c.inputCoverageClient.GetUniqueMethods(filters)...)

		methods, err := c.reflectionClient.GetServiceMethods(logicalService)
		if err != nil {
			return ServiceCoverage{}, err
		}

		reflectionMethods = append(reflectionMethods, methods...)
	}

	return ServiceCoverage{
		TotalCoverage: getCoveragePercent(reflectionMethods, resultsMethods),
	}, nil
}
