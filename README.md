Improved Bogosort
=======
Improved Bogosort is a sorting algorithm generator using state of the art technology,
including machine learning and genetic evolution. Clearly, the machines can write
algorithms better than we can, so I decided to have them sort numbers for me.

Usage
=======
* go run *.go -cmd=generate -programs=<program_file> -num_programs=100
* go run *.go -cmd=evolve -programs=<program_file> -iterations=100
* go run *.go -cmd=test -programs=<program_file>

TODO
=======
[x] Implement all instructions
[x] Randomly generate programs
[x] Run programs on random permuation of data
[x] Rank programs by how well sorted the data is
[x] Save programs after every run
