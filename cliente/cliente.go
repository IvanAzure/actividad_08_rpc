package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"runtime"
)

var clear map[string]func()
type StudentDataRecive struct {
	Name    string
	Subject string
	Score   float64
}

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func client() {
	c, err := rpc.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	var continuar bool = true
	for continuar {
		switch mainMenu() {
		case 1:
			var name string
			var subject string
			var score float64
			var result string
			var ph string
			CallClear()
			fmt.Println("[Agregar la calificacion de un alumno por materia]\n")
			fmt.Print("Alumno:")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				name = scanner.Text()
			} else {
				fmt.Println("Error")
				os.Exit(3)
			}
			fmt.Print("Materia:")
			if scanner.Scan() {
				subject = scanner.Text()
			} else {
				fmt.Println("Error")
				os.Exit(3)
			}
			fmt.Print("Calificacion:")
			fmt.Scanln(&score)

			data := StudentDataRecive{
				Name:    name,
				Subject: subject,
				Score:   score,
			}

			err = c.Call("Server.AddStudentData", data, &result)
			if err != nil {
				fmt.Println(err)

			} else {
				fmt.Println("\n")
				fmt.Println(result)
				fmt.Println("\n")
			}
			fmt.Scanln(&ph)
			break
		case 2:
			var name string
			var result string
			var ph string

			CallClear()
			fmt.Println("[Obtener el promedio del alumno]\n")
			fmt.Print("Alumno:")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				name = scanner.Text()
			} else {
				fmt.Println("Error")
				os.Exit(3)
			}
			err = c.Call("Server.GetStudentAverage", name, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("\n")
				fmt.Println(result)
				fmt.Println("\n")
			}
			fmt.Scanln(&ph)
			break
		case 3:
			var petition = "petition"
			var result string
			var ph string
			err = c.Call("Server.GetGeneralAverageByStudents", petition, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("\n")
				fmt.Println(result)
				fmt.Println("\n")
			}
			fmt.Scanln(&ph)
			break
		case 4:
			CallClear()
			var subject string
			var result string
			var ph string
			fmt.Println("[Obtener el promedio por materia]\n")
			fmt.Print("Materia:")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				subject = scanner.Text()
			} else {
				fmt.Println("Error")
				os.Exit(3)
			}

			err = c.Call("Server.GetAverageBySubject", subject, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("\n")
				fmt.Println(result)
				fmt.Println("\n")
			}
			fmt.Scanln(&ph)
			break
		case 5:
			continuar = false
			break
		default:
			break
		}
	}
}

func mainMenu() int64 {
	var opt int64
	CallClear();
	fmt.Println("[Capturador de calificaciones]")
	fmt.Println("1- Agregar calificacion de un alumno por materia")
	fmt.Println("2- Obtener el promedio del alumno")
	fmt.Println("3- Obtener el promedio de todos los alumnos")
	fmt.Println("4- Obtener el promedio por materia")
	fmt.Println("5- Salir")
	fmt.Print("Opcion:")
	fmt.Scanln(&opt)

	return opt
}

func main() {
	client()
}