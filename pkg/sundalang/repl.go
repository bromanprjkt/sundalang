package sundalang

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

const PROMPT = "sundalang>> "

func StartREPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := NewEnvironment()

	printHeader(out)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		switch line {
		case "kaluar", "exit", "quit":
			fmt.Fprintf(out, "Hatur nuhun parantos nganggo SundaLang!\n")
			return
		case "bantuan", "help":
			printHelp(out)
			continue
		case "bersih", "cls", "clear":
			clearScreen()
			continue
		case "":
			continue
		}

		l := NewLexer(line)
		p := NewParser(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := Eval(program, env)
		if evaluated != nil {
			if evaluated.Type() != NULL_OBJ {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}
		}
	}
}

func printHeader(out io.Writer) {
	fmt.Fprintf(out, "Wilujeng Sumping di SundaLang v1.0.1 (REPL Mode)\n")
	fmt.Fprintf(out, "Ketik 'bantuan' pikeun ningali daptar parentah.\n")
	fmt.Fprintf(out, "Ketik 'kaluar' pikeun udahan.\n")
}

func printHelp(out io.Writer) {
	helpText := `
=== Daptar Parentah REPL ===
  kaluar   : Keluar ti program (atawa ketik 'exit')
  bersih   : Ngabersihkeun layar terminal (atawa ketik 'cls')
  bantuan  : Nampilkeun pesen ieu

=== Conto Syntax ===
  Variabel : tanda x = 10
  Output   : cetakkeun(x)
  Input    : tanda ngaran = tanyakeun("Saha?")
  Kondisi  : lamun x > 5 { cetakkeun("Ageung") }
  Looping  : kedap x > 0 { tanda x = x - 1 }
`
	fmt.Fprint(out, helpText)
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "  Aya nu salah yeuh (Parser Errors):\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}