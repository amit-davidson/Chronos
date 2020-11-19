package output

import (
	"fmt"
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/ssaUtils"
	"github.com/amit-davidson/Chronos/utils"
	"golang.org/x/tools/go/ssa"
	"strings"
	"unicode"
)

const (
	spacePrefixCount = 8
)

func GenerateError(conflictingGAs [][]*domain.GuardedAccess, prog *ssa.Program) error {
	foundDataRaces := utils.NewDoubleKeyMap() // To avoid reporting on the same pair of positions more then once. Can happen if for the same place we read and then write.
	truncatedonflictingGAs := make([][]*domain.GuardedAccess, 0)
	for _, conflict := range conflictingGAs {
		isExist := foundDataRaces.IsExist(conflict[0].Pos, conflict[1].Pos)
		if !isExist {
			foundDataRaces.Add(conflict[0].Pos, conflict[1].Pos)
			truncatedonflictingGAs = append(truncatedonflictingGAs, conflict)
		}
	}

	messages := make([]string, 0)
	for _, conflict := range truncatedonflictingGAs {
		label, err := getMessage(conflict[0], conflict[1], prog)
		if err != nil {
			return err
		}
		messages = append(messages, label)

		if len(messages) == 0 {
			print("No data races found\n")
		}
		print(messages[0])
		for _, message := range messages[1:] {
			print("=========================\n")
			print(message)
		}
	}
	return nil
}

func getMessage(guardedAccessA, guardedAccessB *domain.GuardedAccess, prog *ssa.Program) (string, error) {
	message := "Potential race condition:\n"
	messageA, err := getMessageByLine(guardedAccessA, prog)
	if err != nil {
		return "", err
	}
	messageB, err := getMessageByLine(guardedAccessB, prog)
	if err != nil {
		return "", err
	}
	message += fmt.Sprintf(" %s:\n%s\n \n %s:\n%s \n", "Access1", messageA, "Access2", messageB)
	return message, nil
}

func getMessageByLine(guardedAccessA *domain.GuardedAccess, prog *ssa.Program) (string, error) {
	message := ""
	posA := prog.Fset.Position(guardedAccessA.Pos)
	lineA, err := utils.ReadLineByNumber(posA.Filename, posA.Line)
	trimmedA := strings.TrimLeftFunc(lineA, unicode.IsSpace)
	message += strings.Repeat(" ", spacePrefixCount) + trimmedA

	removedSpaces := len(lineA) - len(trimmedA)
	posToAddArrow := posA.Column - removedSpaces
	message += "\n" + strings.Repeat(" ", posToAddArrow+spacePrefixCount-1) + "^" + "\n"
	stackA := ssaUtils.GetStackTrace(prog, guardedAccessA)
	if err != nil {
		return "", err
	}
	stackA += posA.String()
	message += stackA
	return message, nil
}
