go build -o ../output/GoNES \
    ../internal/Main.go \
    ../internal/Bus.go \
    ../internal/C6502.go \
    ../internal/P2C02.go \
    ../internal/DataTypes.go \
    ../internal/Utils.go \
    ../internal/Cartridge.go \
    ../internal/Mapper.go \
    ../internal/Debug.go
