
res = 'texts/results.txt'
btc = 'texts/btc.txt'

target = res
# target = btc

def main():
    ll:set = None
    with open(target,"r", encoding="utf8") as f:
        ll = set(f.read().splitlines())

    with open(target,"w", encoding="utf8") as f:
        [f.write(i + '\n') for i in ll]

if __name__ == "__main__":
    main()