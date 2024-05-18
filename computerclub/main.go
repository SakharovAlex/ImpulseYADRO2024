package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const timeLayout = "15:04"

type Event struct {
	Time      time.Time
	ID        int
	EventBody []string
}

type Table struct {
	OccupiedBy   string
	StartTime    time.Time
	Revenue      int
	OccupiedTime int
}

type ClientInfo struct {
	StartTime   time.Time
	PlaceNumber int
}

type Club struct {
	OpenTime  time.Time
	CloseTime time.Time
	HourPrice int
	Events    []Event
	Clients   map[string]int
	Queue     []string
	Tables    []Table
}

func parseNumber(scanner *bufio.Scanner) (int, error) {
	scanner.Scan()
	numTables, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, fmt.Errorf(scanner.Text())
	}
	return numTables, nil
}

func parseOpenAndCloseTimes(scanner *bufio.Scanner) (time.Time, time.Time, error) {
	scanner.Scan()
	times := strings.Split(scanner.Text(), " ")
	openTime, err := time.Parse(timeLayout, times[0])
	if err != nil {
		return time.Time{}, time.Time{},
			fmt.Errorf(scanner.Text())
	}
	closeTime, err := time.Parse(timeLayout, times[1])
	if err != nil {
		return time.Time{}, time.Time{},
			fmt.Errorf(scanner.Text())
	}
	return openTime, closeTime, nil
}

func isValidName(s string) bool {
	re := regexp.MustCompile(`^[a-z0-9_-]+$`)
	return re.MatchString(s)
}

func parseEvent(scanner *bufio.Scanner) (Event, error) {
	parts := strings.Split(scanner.Text(), " ")
	if len(parts) < 3 {
		return Event{}, fmt.Errorf(scanner.Text())
	}
	eventTime, err := time.Parse(timeLayout, parts[0])
	if err != nil {
		return Event{}, err
	}
	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return Event{}, err
	}
	if !isValidName(parts[2]) {
		return Event{}, fmt.Errorf(scanner.Text())
	}
	return Event{
		Time:      eventTime,
		ID:        eventID,
		EventBody: parts[2:],
	}, nil
}

func parseInput() (Club, error) {
	if len(os.Args) < 2 {
		return Club{}, fmt.Errorf("not enough arguments")
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		return Club{}, err
	}
	defer file.Close()
	club := Club{}
	scanner := bufio.NewScanner(file)
	numTables, err := parseNumber(scanner)
	if err != nil {
		return Club{}, err
	}
	club.Tables = make([]Table, numTables)
	club.OpenTime, club.CloseTime, err = parseOpenAndCloseTimes(scanner)
	if err != nil {
		return Club{}, err
	}
	club.HourPrice, err = parseNumber(scanner)
	if err != nil {
		return Club{}, err
	}
	var lastEventTime time.Time
	for scanner.Scan() {
		event, err := parseEvent(scanner)
		if err != nil {
			return Club{}, fmt.Errorf(scanner.Text())
		}
		if len(club.Events) != 0 {
			if event.Time.Before(lastEventTime) {
				return Club{}, fmt.Errorf(scanner.Text())
			}
		}
		club.Events = append(club.Events, event)
		lastEventTime = event.Time
	}
	if err := scanner.Err(); err != nil {
		return Club{}, fmt.Errorf("%v", err)
	}
	club.Clients = make(map[string]int)
	club.Queue = make([]string, 0)
	return club, nil
}

func (club *Club) handleEvent(event Event) string {
	switch event.ID {
	case 1:
		return club.clientCome(event)
	case 2, 12:
		return club.clientTakesSeat(event)
	case 3:
		return club.clientWaits(event)
	case 4:
		return club.clientLeaves(event)
	default:
		return fmt.Sprintf("%s 13 IncorrectEventID", event.Time.Format(timeLayout))
	}
}

func (club *Club) clientCome(event Event) string {
	clientName := event.EventBody[0]
	fmt.Printf("%s %d %s\n", event.Time.Format(timeLayout), event.ID, clientName)
	if _, exists := club.Clients[clientName]; exists {
		return fmt.Sprintf("%s 13 YouShallNotPass", event.Time.Format(timeLayout))
	}
	if event.Time.Before(club.OpenTime) || event.Time.After(club.CloseTime) {
		return fmt.Sprintf("%s 13 NotOpenYet", event.Time.Format(timeLayout))
	}
	club.Clients[clientName] = 0
	return ""
}

func (club *Club) clientTakesSeat(event Event) string {
	// если клиент сидит за столом, то он может сменить его
	clientName := event.EventBody[0]
	tableNumber, err := strconv.Atoi(event.EventBody[1])
	if err != nil {
		return fmt.Sprintf("%s 13 IncorrectTableNumber", event.Time.Format(timeLayout))
	}
	fmt.Printf("%s %d %s %d\n", event.Time.Format(timeLayout), event.ID, clientName, tableNumber)
	if tableNumber < 1 || tableNumber > len(club.Tables) {
		return fmt.Sprintf("%s 13 IncorrectTableNumber", event.Time.Format(timeLayout))
	}
	newTable := &club.Tables[tableNumber-1]
	if place, exists := club.Clients[clientName]; !exists {
		return fmt.Sprintf("%s 13 ClientUnknown", event.Time.Format(timeLayout))
	} else if newTable.OccupiedBy != "" {
		return fmt.Sprintf("%s 13 PlaceIsBusy", event.Time.Format(timeLayout))
	} else { // Клиент существует и стол свободный
		// Для того чтобы уйти с того места, где сидели
		if place != 0 {
			club.clientLeavesTable(place, event.Time)
		}
		club.Clients[clientName] = tableNumber
		newTable.OccupiedBy = clientName
		newTable.StartTime = event.Time
	}
	return ""
}

func (club *Club) clientWaits(event Event) string {
	clientName := event.EventBody[0]
	fmt.Printf("%s %d %s\n", event.Time.Format(timeLayout), event.ID, clientName)
	if _, exists := club.Clients[clientName]; !exists {
		return fmt.Sprintf("%s 13 ClientUnknown", event.Time.Format(timeLayout))
	}
	for _, table := range club.Tables {
		if table.OccupiedBy == "" {
			return fmt.Sprintf("%s 13 ICanWaitNoLonger!", event.Time.Format(timeLayout))
		}
	}
	if len(club.Queue) >= len(club.Tables) {
		return fmt.Sprintf("%s 11 %s", event.Time.Format(timeLayout), clientName)
	}
	club.Queue = append(club.Queue, clientName)
	return ""
}

func (club *Club) clientLeaves(event Event) string {
	clientName := event.EventBody[0]
	fmt.Printf("%s %d %s\n", event.Time.Format(timeLayout), event.ID, clientName)
	place, exists := club.Clients[clientName]
	if !exists {
		return fmt.Sprintf("%s 13 ClientUnknown", event.Time.Format(timeLayout))
	}
	if place != 0 {
		club.clientLeavesTable(place, event.Time)
	}
	delete(club.Clients, clientName)
	return ""
}

func (club *Club) clientLeavesTable(tableNumber int, leavingTime time.Time) {
	table := &club.Tables[tableNumber-1]
	duration := leavingTime.Sub(table.StartTime)
	hours := int(duration.Hours())
	if duration.Minutes() > float64(hours*60) {
		hours++
	}
	table.Revenue += hours * club.HourPrice
	table.OccupiedTime += int(duration.Minutes())
	table.OccupiedBy = ""
	if len(club.Queue) > 0 && leavingTime != club.CloseTime { // МБ второе условие нужно убрать
		nextClient := club.Queue[0]
		club.Queue = club.Queue[1:]
		club.handleEvent(Event{
			Time:      leavingTime,
			ID:        12,
			EventBody: []string{nextClient, strconv.Itoa(tableNumber)},
		}) // Возможно стоит печатать возращаемое значение
	}
}

func (club *Club) closureProcessing() {
	sortedClients := make([]string, 0)
	for client, _ := range club.Clients {
		sortedClients = append(sortedClients, client)
	}
	sort.Strings(sortedClients)
	for _, client := range sortedClients {
		place := club.Clients[client]
		fmt.Printf("%s 11 %s\n", club.CloseTime.Format(timeLayout), client)
		if place != 0 {
			club.clientLeavesTable(place, club.CloseTime)
		}
	}
}

func (club *Club) processEvents() {
	for _, event := range club.Events {
		if result := club.handleEvent(event); result != "" {
			fmt.Println(result)
		}
	}
	club.closureProcessing()
}

func (club *Club) workReport() {
	fmt.Println(club.OpenTime.Format(timeLayout))
	club.processEvents()
	fmt.Println(club.CloseTime.Format(timeLayout))
	for i, table := range club.Tables {
		tableHours, tableMinutes := table.OccupiedTime/60, table.OccupiedTime%60
		fmt.Printf("%d %d %02d:%02d\n", i+1, table.Revenue, tableHours, tableMinutes)
	}
}

func main() {
	club, err := parseInput()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
		return
	}
	club.workReport()
}
