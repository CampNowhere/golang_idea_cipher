# IDEA Cipher For Go

This is code to allow use of the IDEA encryption algorithm in golang.

It is a block cipher, with a 64 bit block size and a 128 bit key size. 

It implements Block interface found in the standard cipher package so it can drop in to any existing code that uses block ciphers. 

## Note on Security

As of 2013, the only known faster-than-bruteforce attack is not computationally feasible. However, the small block size of IDEA does not make it preferential for new applications. From a pragmatic standpoint, I woud say that this cipher is *probably* secure, but from a security perspective, I can't in good conscience recommend that this cipher be used for any purposes but legacy compatibility. Also, the function called within each round relies heavily upon looping, which may leak information about the state of the cipher. Also (keep in mind that this is my own editorial opinion), I believe the key schedule is a bit weak for this algorithm. Generation of each round key involves simply rotating bits, with no whitening or s-box permutation. Therefore, a key with a lot of zeros will generate a key schedule with a lot of zeros. 

## Note on Performance

The complexity of the round function makes this cipher a bit slower than most. The necessity of calling the imul function six times per round (and with no clear way of inlining functions or doing macro expansion in Go) adds considerable overhead, and the imul function itself contains a loop.
