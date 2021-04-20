with open('file-1.txt', 'w') as f:
    for i in range(10000000):
        f.write(f'{i:09}\n')
    