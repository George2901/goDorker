

target = 'texts/results.txt'
# target = btc

def main():
    ll:list = None
    with open(target,"r", encoding="utf8") as f:
        ll = f.read().splitlines()

    for i,j in enumerate(ll):
        if j.startswith('http'):pass
        elif j.startswith('https'):pass
        else:ll.pop(i)
        

    with open(target,"w", encoding="utf8") as f:
        [f.write(i + '\n') for i in ll]

if __name__ == "__main__":
    main()