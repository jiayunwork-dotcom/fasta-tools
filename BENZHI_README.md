核酸序列 Go 命令行：子命令 gc / rc / kmer 读 FASTA 算 GC 含量、反向互补和 k-mer 谱，fasta-tools serve 把这三项变成 HTTP。

# fasta-tools

Bioinformatics sequence analysis service providing GC content calculation,
reverse complement generation, k-mer frequency counting, and FASTA parsing
via HTTP API and CLI.

## Build / Run / Test

```bash
go build -o fasta-tools .
./fasta-tools serve -addr :8080
go test ./...
```

## Evaluation Image

Evaluation-specific files (do not overwrite project Dockerfile/README):

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md` (this file)

Build and verify in container:

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
