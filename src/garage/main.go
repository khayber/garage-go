package main

func main() {
    setup()
    go monitor()
    go tgbot()
    rest()
    cleanup()
}

