parts

Partition and Join large files for easier and more reliable transmission

Usage:

  - Partition

    ```parts -f filename.ext -s 40000000```

    Partitions a file into 40 million byte chunks

  - Join

    ```part -f filename.ext -m join```

    Joins all file parts into filename.ext where there exists files with the pattern filename.ext.0 ... filename.ext.N