package main

import (
	"fmt"
	"testing"
)

func TestParseWithoutArguments(t *testing.T) {
	command := ParseCommand("/chat_id")
	if command.name != "chat_id" {
		t.Error(fmt.Sprintf("expected chat_id, got %s", command.name))
	}
}

func TestNotCommand(t *testing.T) {
	command := ParseCommand("ololol")
	if command.name != "" {
		t.Error(fmt.Sprintf("expected '', got %s", command.name))
	}
}

func TestParseWithSpaces(t *testing.T) {
	command := ParseCommand("/chat_id                        ")
	if command.name != "chat_id" {
		t.Error(fmt.Sprintf("expected chat_id, got %s", command.name))
	}
}

func TestParseWithArg(t *testing.T) {
	command := ParseCommand("/stats LUV_nellrun")
	if command.name != "stats" {
		t.Error(fmt.Sprintf("expected stats, got %s", command.name))
	}

	if len(command.args) != 1 {
		t.Error(fmt.Sprintf("unexpected number of args %d", len(command.args)))
	}

	if command.args[0] != "LUV_nellrun" {
		t.Error(fmt.Sprintf("unexpected arg '%s'", command.args[0]))
	}
}

func TestParseTrimSpaces(t *testing.T) {
	command := ParseCommand("/stats                    LUV_nellrun                           ")
	if command.name != "stats" {
		t.Error(fmt.Sprintf("expected stats, got %s", command.name))
	}

	if len(command.args) != 1 {
		t.Error(fmt.Sprintf("unexpected number of args %d", len(command.args)))
	}

	if command.args[0] != "LUV_nellrun" {
		t.Error(fmt.Sprintf("unexpected arg '%s'", command.args[0]))
	}
}

func TestParseWithArgs(t *testing.T) {
	command := ParseCommand("/stats    LUV_nellrun     lifeline       ")
	if command.name != "stats" {
		t.Error(fmt.Sprintf("expected stats, got %s", command.name))
	}

	if len(command.args) != 2 {
		t.Error(fmt.Sprintf("unexpected number of args %d", len(command.args)))
	}

	if command.args[0] != "LUV_nellrun" || command.args[1] != "lifeline" {
		t.Error(fmt.Sprintf("unexpected args '%s', '%s'", command.args[0], command.args[1]))
	}
}
