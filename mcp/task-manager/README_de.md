# Task Manager MCP Server

Ein hochentwickelter State-Management-Server, der es Agenten ermöglicht, eine persistente Checkliste von Zielen und Architekturentscheidungen zu führen.

## Kernkonzepte

### Plan-Modus
Ein spezieller Zustand, in dem der Agent auf "Read-Only"-Operationen beschränkt ist. In dieser Phase exploriert der Agent die Codebasis, entwirft eine Implementierungsstrategie und füllt die Aufgabenliste. Der Agent muss diesen Plan dem Benutzer präsentieren und eine Genehmigung einholen, bevor er den Plan-Modus verlässt, um mit der Implementierung zu beginnen.

## Tools

### 1. `mlc_task_create`
Erstellt eine neue strukturierte Aufgabe in der Checkliste.
- **Anwendungsfall**: Verfolgung des Fortschritts bei komplexen mehrstufigen Refactorings.

### 2. `mlc_task_update`
Aktualisiert den Status einer Aufgabe (`pending`, `in_progress`, `completed`, `deleted`) oder verwaltet Abhängigkeiten.

### 3. `mlc_task_list`
Zeigt alle Aufgaben der aktuellen Sitzung an.

### 4. `mlc_task_get`
Ruft detaillierte Informationen zu einer bestimmten Aufgabe ab.

### 5. `mlc_enter_plan_mode`
Versetzt den Agenten in den Plan-Modus. Sollte proaktiv bei jeder nicht-trivialen Aufgabe genutzt werden.

### 6. `mlc_exit_plan_mode`
Versetzt den Agenten zurück in den Implementierungs-Modus, nachdem der Benutzer dem Plan zugestimmt hat.

## Installation

Wird als Teil des Hauptprojekts gebaut:
```bash
task build
```
