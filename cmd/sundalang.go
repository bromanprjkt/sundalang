package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sundalang/pkg/sundalang"
)

func runCode(code string) {
	l := sundalang.NewLexer(code)
	p := sundalang.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Printf("Parser Errors:\n")
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
		}
		return
	}

	env := sundalang.NewEnvironment()
	evalResult := sundalang.Eval(program, env)
	if evalResult != nil {
		if evalResult.Type() == sundalang.ERROR_OBJ {
			fmt.Println(evalResult.Inspect())
		}
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Gagal maca file: %s\n", err)
		return
	}
	runCode(string(content))
}

func showHelp() {
	fmt.Println("SundaLang - Bahasa Pemrograman Sunda Pandeglang")
	fmt.Println()
	fmt.Println("Cara make:")
	fmt.Println("  sundalang <file.sl> [arg...]  Jalankeun file SundaLang")
	fmt.Println("  sundalang run <file.sl>       Jalankeun file SundaLang")
	fmt.Println("  sundalang repl                Buka mode interaktif (REPL)")
	fmt.Println("  sundalang -e '<code>'         Jalankeun kode saharita (inline eval)")
	fmt.Println("  sundalang                     Buka Launcher / Menu Utama")
	fmt.Println("  sundalang --version           Tempo versi SundaLang")
	fmt.Println("  sundalang install             Install SundaLang ka sistem")
	fmt.Println("  sundalang uninstall           Uninstall SundaLang tina sistem")
	fmt.Println("  sundalang --help              Tempo pitulung ieu")
	fmt.Println()
}

func install() {
	setupInstallerWindow()
	defer cleanupInstallerWindow()
	printBanner()
	printHeader("SELF-INSTALLER")

	exePath, err := os.Executable()
	if err != nil {
		printError(fmt.Sprintf("Gagal manggihan lokasi executable: %s", err))
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		printError(fmt.Sprintf("Gagal manggihan home directory: %s", err))
		return
	}

	installDir := filepath.Join(homeDir, ".sundalang", "bin")

	printInfo(fmt.Sprintf("Nyieun folder instalasi: %s", installDir))
	if err := os.MkdirAll(installDir, 0755); err != nil {
		printError(fmt.Sprintf("Gagal nyieun folder: %s", err))
		return
	}

	targetName := "sundalang"
	if runtime.GOOS == "windows" {
		targetName = "sundalang.exe"
	}
	targetPath := filepath.Join(installDir, targetName)

	printInfo(fmt.Sprintf("Nyalin executable ka: %s", targetPath))
	if err := copyFile(exePath, targetPath); err != nil {
		printError(fmt.Sprintf("Gagal nyalin file: %s", err))
		return
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(targetPath, 0755); err != nil {
			printError(fmt.Sprintf("Gagal set permissions: %s", err))
			return
		}
	}

	printSuccess("File geus dicopy")
	fmt.Println()

	printInfo("Ngatur PATH...")
	if runtime.GOOS == "windows" {
		if err := addToWindowsPath(installDir); err != nil {
			printWarning(fmt.Sprintf("Gagal otomatis nambahkeun ka PATH: %s", err))
			fmt.Println()
			printLine(ColorRed, "  Tambahkeun sacara manual:")
			printLine(ColorRed, "  Buka System Properties > Environment Variables")
			printLine(ColorRed, fmt.Sprintf("  Tambahkeun '%s' ka User PATH", installDir))
		} else {
			printSuccess("PATH geus diupdate")
		}
	} else {
		shellRC := filepath.Join(homeDir, ".bashrc")
		if _, err := os.Stat(filepath.Join(homeDir, ".zshrc")); err == nil {
			shellRC = filepath.Join(homeDir, ".zshrc")
		}

		content, _ := os.ReadFile(shellRC)
		if !strings.Contains(string(content), ".sundalang/bin") {
			f, err := os.OpenFile(shellRC, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString("\n# SundaLang\n")
				f.WriteString("export PATH=\"$HOME/.sundalang/bin:$PATH\"\n")
				f.Close()
				printSuccess(fmt.Sprintf("PATH geus ditambahkeun ka %s", shellRC))
				printInfo(fmt.Sprintf("Jalankeun: source %s", shellRC))
			}
		} else {
			printSuccess("PATH geus aya")
		}
	}

	fmt.Println()
	printSuccessBox("INSTALASI BERHASIL!")
	fmt.Println()
	printSeparator()
	fmt.Println()

	if runtime.GOOS == "windows" {
		printWarning("PENTING: Tutup tur buka deui terminal pikeun nerapkeun PATH")
		fmt.Println()
	}

	printInfo("jalankeun jeng cek versi: " + colorize(ColorRed+ColorBold, "sundalang --version"))
	fmt.Println()
	printFooter()
	printPrompt(colorize(ColorWhite, "Pencet Enter pikeun nutup... "))
	fmt.Scanln()
}

func uninstall() {
	setupInstallerWindow()
	defer cleanupInstallerWindow()
	printBanner()
	printHeader("SELF-UNINSTALLER")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		printError(fmt.Sprintf("Gagal manggihan home directory: %s", err))
		return
	}

	installDir := filepath.Join(homeDir, ".sundalang")

	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		printInfo("SundaLang teu kapanggih. Meureun geus diuninstall.")
		fmt.Println()
		printPrompt(colorize(ColorGray, "Pencet Enter pikeun nutup..."))
		fmt.Scanln()
		return
	}

	printInfo(fmt.Sprintf("Kapanggih instalasi di: %s", installDir))
	fmt.Println()
	printFooter()
	printPrompt(colorize(ColorWhite, "Rek diuninstall SundaLang? (y/n): "))

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		printInfo("Uninstall dibatalkeun.")
		fmt.Println()
		printPrompt(colorize(ColorGray, "Pencet Enter pikeun nutup..."))
		fmt.Scanln()
		return
	}

	fmt.Println()
	printInfo("Ngahapus file...")

	if runtime.GOOS == "windows" {
		binPath := filepath.Join(installDir, "bin")
		removeFromWindowsPath(binPath)
	}

	err = os.RemoveAll(installDir)
	if err != nil {
		if runtime.GOOS == "windows" {
			printWarning("File sedang dipake, ngagunakeun delayed delete...")
			createWindowsDelayedDelete(installDir)
			printSuccess("Uninstaller dijadwalkeun. File bakal dihapus sanggeus program nutup.")
		} else {
			printError(fmt.Sprintf("Gagal ngahapus folder: %s", err))
			fmt.Println()
			printFooter()
			printPrompt(colorize(ColorWhite, "Pencet Enter pikeun nutup... "))
			fmt.Scanln()
			return
		}
	} else {
		printSuccess("File geus dihapus")
	}

	fmt.Println()
	printSeparator()
	printSuccessBox("UNINSTALL BERHASIL!")
	fmt.Println()
	printSeparator()
	fmt.Println()

	printLine(ColorWhite, "Hatur nuhun geus make SundaLang!")
	fmt.Println()
	printInfo("Tutup terminal tur buka deui pikeun nerapkeun parobahan PATH")
	fmt.Println()
	printFooter()
	printPrompt(colorize(ColorWhite, "Pencet Enter pikeun nutup... "))
	fmt.Scanln()
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0755)
}

func showMainMenu() {
	setupInstallerWindow()
	defer cleanupInstallerWindow()
	printBanner()
	printHeader("INSTALLER & LAUNCHER")
	fmt.Println()

	printMenuOption(1, "Install SundaLang ka sistem")
	printMenuOption(2, "Uninstall SundaLang")
	printMenuOption(3, "Jalankeun file .sl")
	printMenuOption(4, "Buka REPL mode")
	printMenuOption(5, "Keluar")

	fmt.Println()
	printSeparator()
	fmt.Println()
	fmt.Println()
	printFooter()
	printPrompt(colorize(ColorWhite, "Pilihan (1-5): "))

	var choice string
	fmt.Scanln(&choice)
	fmt.Println()

	switch choice {
	case "1":
		install()
	case "2":
		uninstall()
	case "3":
		printFooter()
		printPrompt(colorize(ColorWhite, "Asupkeun path file .sl: "))
		var filePath string
		fmt.Scanln(&filePath)
		if filePath != "" {
			fmt.Println()
			runFile(filePath)
		}
		fmt.Println()
		fmt.Println()
		printPrompt(colorize(ColorGray, "Pencet Enter pikeun nutup..."))
		fmt.Scanln()
	case "4":
		fmt.Println()
		sundalang.StartREPL(os.Stdin, os.Stdout)
	case "5":
		fmt.Println()
		printLine(ColorWhite, "Hatur nuhun!")
		return
	default:
		printError("Pilihan teu valid")
		fmt.Println()
		fmt.Println()
		printPrompt(colorize(ColorGray, "Pencet Enter pikeun nutup..."))
		fmt.Scanln()
	}
}

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]

		switch arg {
		case "-v", "--version":
			fmt.Println("SundaLang v1.0.5")
			fmt.Println("Bahasa Pemrograman Sunda Pandeglang")
			return
		case "--help", "-h", "help":
			showHelp()
			return
		case "install":
			install()
			return
		case "uninstall":
			uninstall()
			return
		case "repl":
			sundalang.StartREPL(os.Stdin, os.Stdout)
			return
		case "-e", "--eval":
			if len(os.Args) < 3 {
				fmt.Println("Kudu aya kode nu dievaluasi. Conto: sundalang -e 'cetakkeun(\"Halo\")'")
				return
			}
			evalArgs := os.Args[2:]
			dashDashIdx := -1
			for i, a := range evalArgs {
				if a == "--" {
					dashDashIdx = i
					break
				}
			}

			var code string
			if dashDashIdx != -1 {
				code = strings.Join(evalArgs[:dashDashIdx], " ")
				sundalang.SetArgs(evalArgs[dashDashIdx+1:])
			} else {
				code = strings.Join(evalArgs, " ")
				sundalang.SetArgs([]string{})
			}
			runCode(code)
			return
		case "run":
			if len(os.Args) < 3 {
				fmt.Println("Kudu aya ngaran file .sl nu rek dijalankeun. Conto: sundalang run contoh.sl")
				return
			}
			if len(os.Args) > 3 {
				sundalang.SetArgs(os.Args[3:])
			} else {
				sundalang.SetArgs([]string{})
			}
			runFile(os.Args[2])
			return
		default:
			if len(os.Args) > 2 {
				sundalang.SetArgs(os.Args[2:])
			} else {
				sundalang.SetArgs([]string{})
			}
			runFile(arg)
			return
		}
	}

	showMainMenu()
}
