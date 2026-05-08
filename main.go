package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
)

type Game struct {
	PlayerName   string
	MonsterName  string
	PlayerHP     int
	MonsterHP    int
	MaxPlayerHP  int
	MaxMonsterHP int
	Message      string
	LastAction   string
	GameOver     bool
}

var game Game

var page = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	newGame()

	http.HandleFunc("/", showGame)
	http.HandleFunc("/attack", attack)
	http.HandleFunc("/heal", heal)
	http.HandleFunc("/super", superAttack)
	http.HandleFunc("/restart", restart)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Go Dungeon запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func newGame() {
	// Дети могут менять эти значения.
	playerName := "Рыцарь Go"
	monsterName := "Багозавр"
	playerHP := 100
	monsterHP := 1000

	game = Game{
		PlayerName:   playerName,
		MonsterName:  monsterName,
		PlayerHP:     playerHP,
		MonsterHP:    monsterHP,
		MaxPlayerHP:  playerHP,
		MaxMonsterHP: monsterHP,
		Message:      "Ты вошел в подземелье. Багозавр уже ждет!",
		LastAction:   "start",
		GameOver:     false,
	}
}

func (g Game) PlayerPercent() int {
	return percent(g.PlayerHP, g.MaxPlayerHP)
}

func (g Game) MonsterPercent() int {
	return percent(g.MonsterHP, g.MaxMonsterHP)
}

func percent(value int, max int) int {
	if max <= 0 || value <= 0 {
		return 0
	}
	if value >= max {
		return 100
	}
	return value * 100 / max
}

func showGame(w http.ResponseWriter, r *http.Request) {
	page.Execute(w, game)
}

func attack(w http.ResponseWriter, r *http.Request) {
	if game.GameOver {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	damage := rand.IntN(20) + 5
	message := fmt.Sprintf("Ты ударил монстра на %d урона.", damage)

	if rand.IntN(100) < 20 {
		damage *= 2
		message = fmt.Sprintf("Критический удар! Ты нанес %d урона.", damage)
	}

	game.MonsterHP -= damage
	game.LastAction = "attack"
	monsterTurn(message)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func heal(w http.ResponseWriter, r *http.Request) {
	if game.GameOver {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	heal := rand.IntN(15) + 5
	game.PlayerHP += heal

	if game.PlayerHP > game.MaxPlayerHP {
		game.PlayerHP = game.MaxPlayerHP
	}

	game.LastAction = "heal"
	message := fmt.Sprintf("Ты выпил зелье и восстановил %d HP.", heal)
	monsterTurn(message)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func superAttack(w http.ResponseWriter, r *http.Request) {
	if game.GameOver {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	game.LastAction = "super"

	if rand.IntN(100) < 50 {
		damage := rand.IntN(35) + 15
		game.MonsterHP -= damage
		monsterTurn(fmt.Sprintf("Суперудар попал! Багозавр получил %d урона.", damage))
	} else {
		monsterTurn("Суперудар промахнулся!")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func monsterTurn(playerMessage string) {
	if game.MonsterHP <= 0 {
		game.MonsterHP = 0
		game.GameOver = true
		game.LastAction = "win"
		game.Message = playerMessage + " Победа! Монстр побежден."
		return
	}

	damage := rand.IntN(12) + 4
	game.PlayerHP -= damage

	if game.PlayerHP <= 0 {
		game.PlayerHP = 0
		game.GameOver = true
		game.LastAction = "lose"
		game.Message = playerMessage + fmt.Sprintf(" %s ударил тебя на %d урона. Ты проиграл.", game.MonsterName, damage)
		return
	}

	game.Message = playerMessage + fmt.Sprintf(" %s ударил тебя на %d урона.", game.MonsterName, damage)
}

func restart(w http.ResponseWriter, r *http.Request) {
	newGame()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
