package reflection

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/samber/lo"
)

func GetFieldNamesFromMessageDescriptor(descriptor *desc.MessageDescriptor) []string {
	fields := descriptor.GetFields()
	return lo.Map(fields, func(item *desc.FieldDescriptor, _ int) string { return item.GetName() })
}
