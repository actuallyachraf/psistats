import random

GAUSS_DIST_MU = 50
GAUSS_DIST_SIGMA = 15

def generate_test_data(filename, count):
    with open(filename,'w') as f:
        cols = f"id_a,stat_a,id_b,stat_b\n"
        f.write(cols)
        for i in range(1,count):
            stat_a = int(abs(random.gauss(GAUSS_DIST_MU,GAUSS_DIST_SIGMA)))
            stat_b = int(abs(random.gauss(GAUSS_DIST_MU,GAUSS_DIST_SIGMA)))
            line = f"{i},{stat_a},{i},{stat_b}" + '\n'
            f.write(line)
    f.close()

if __name__ == '__main__':
    generate_test_data("ads-test.csv",500)