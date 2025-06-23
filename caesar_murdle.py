"""
    Caesar Cipher - http://inventwithpython.com/hacking (BSD Licensed)
    modified by mlg

    For the guess command, use the chi squared statistic to choose the best
    key.

"""
import argparse
from collections import defaultdict
import json
import pyperclip
import string
import sys

MODE_ENCODE = 1
MODE_DECODE = 2
MODE_GUESS = 3
LETTERS_UC = set(string.ascii_uppercase)
LETTERS_LC = set(string.ascii_lowercase)
ALEN = 26  # number of letters in the alphabet
DEFAULT_JSON = 'digrams.json'

english_frequences = [0.08167, 0.01492, 0.02782, 0.04253, 0.12702, 0.02228,
                      0.02015, 0.06094, 0.06966, 0.00153, 0.00772, 0.04025,
                      0.02406, 0.06749, 0.07507, 0.01929, 0.00095, 0.05987,
                      0.06327, 0.09056, 0.02758, 0.00978, 0.02360, 0.00150,
                      0.01974, 0.00074]


def get_fscore(msg, key):
    """
    Count the frequency of the letters and compare that with the expected
    frequency of English.

    key - a number in 0..25
    """
    ord_uca = ord('A')
    ord_ucak = ord_uca - key
    tally = defaultdict(int)
    for symbol in msg:
        symbol = symbol.upper()
        if symbol in LETTERS_UC:
            tally[chr((ord(symbol) - ord_ucak) % ALEN + ord_uca)] += 1
    total = float(sum(tally.values()))
    fscore = 0
    chisquare = 0
    # print(tally)
    for symbol in tally:
        percent = tally[symbol] / total
        english_percent = english_frequences[ord(symbol) - ord_uca]
        expected = english_percent * total
        chisquare += (tally[symbol] - expected) ** 2 / expected
    if args.verbose > 1:
        print(f'{key:2} {chisquare=:5.3f}')
    return chisquare


def guess(msg):
    """
    Iterate over the possible keys and return the one with the best digraph
    frequency score.

    Accumulate the scores into a list and if option --outjson is specified,
    write the list to the named file.
    """
    best_fkey = -1
    min_fscore = sys.maxsize
    fscores = []
    for key in range(ALEN):
        fscore = get_fscore(msg, key)
        fscores.append(fscore)
        if fscore < min_fscore:
            best_fkey = key
            min_fscore = fscore
    return best_fkey, fscores


def translate(msg, key):
    """
    Translate the message based on the given key. If the translation is to
    decode the message, the key will have been transformed from the normal
    encode value to the decode value. For example, if the given key is 2, the
    decode value will be -2 modulo 26 or 24.
    """
    ord_uca = ord('A')
    ord_lca = ord('a')
    translated = []
    # run the encryption/decryption code on each symbol in the message string
    for symbol in msg:
        if symbol in LETTERS_UC:
            # get the encrypted (or decrypted) number for this symbol
            num = ord(symbol) - ord_uca
            num = (num + key) % ALEN
            # add encrypted/decrypted number's symbol at the end of translated
            symbol = chr(num + ord_uca)
        elif symbol in LETTERS_LC:
            num = ord(symbol) - ord_lca
            num = (num + key) % ALEN
            symbol = chr(num + ord_lca)
        translated.append(symbol)
    return ''.join(translated)


def get_ciphertext():
    if args.murdle_num is None:
        return args.file.read()
    infile = args.file
    for line in infile:
        line = line.strip()
        if len(line) == 0 or line[0] == '#':
            continue
        ciphernum = int(line[:2])
        if ciphernum == args.murdle_num:
            return line[2:]
    return None


def get_args():
    parser = argparse.ArgumentParser(description='Caesar Cipher encode/decode')
    parser.add_argument('mode', help='''
        "encode", "decode", or "guess",
        may be abbreviated to the first letter.  If "guess" is selected,
        the program will decode using a digraph frequency table to guess
        the key.''')
    parser.add_argument('file', type=argparse.FileType('r'), help='''
        the input file to be encoded or decoded. The result is printed to
        the system output.''')
    parser.add_argument('-k', '--key', type=int, default=13, help='''
        the offset to be added to a letter for encoding or subtracted from it
        for decoding. The default is 13, in which case encoding and decoding
        are identical. This is "rot13" coding.  If the mode is "guess",
        this parameter is ignored.''')
    parser.add_argument('-c', '--clipboard', action='store_true', help='''
        copy the translated message to the clipboard.''')
    parser.add_argument('-j', '--json', default=DEFAULT_JSON, help='''
        if the mode is "guess", this argument specifies the json file
        containing the digraph frequency table. The default is "{}".'''.format(
            DEFAULT_JSON))
    parser.add_argument('-m', '--murdle_num', type=int, help='''
        If present, this is the number of the cipher corresponding to the first two
        columns in a row in the input file. The remainder of the row is the cipher.''')
    parser.add_argument('-o', '--outjson', help='''
        output file name to contain the key scores in json format.  This
        option is only relevant if mode is "guess".''')
    parser.add_argument('-p', '--patrist', action='store_true', help='''
        the cipher is a "patristocrat" with word divisions suppressed. Ignore
        word boundaries. This only makes sense for mode "guess".''')
    parser.add_argument('-u', '--upper', action='store_true', help='''
        convert the text to upper case before processing.''')
    parser.add_argument('-v', '--verbose', type=int, default=1, help='''
        set the verbosity level to 0, 1, or 2. Default = 1''')
    args = parser.parse_args()
    if args.mode[0].lower() == 'e':  # encode
        args.mode = MODE_ENCODE
    elif args.mode[0].lower() == 'd':  # decode
        args.mode = MODE_DECODE
        args.key = ALEN - args.key
    elif args.mode[0].lower() == 'g':  # guess
        args.mode = MODE_GUESS
    else:
        raise ValueError(
            'First argument must be "encode", "decode", or "guess".')
    return args


if __name__ == '__main__':
    assert sys.version_info >= (3, 8)
    args = get_args()
    message = get_ciphertext()
    if message is None:
        print('Invalid cipher number.')
        sys.exit(1)
    args.file.close()  # in case we overwrite it with redirected output
    key = args.key
    if args.upper:
        message = message.upper()
    if args.patrist:
        message = ''.join([c for c in message if c in string.ascii_letters])
    if args.mode == MODE_GUESS:
        key, fscores = guess(message)
        # guess() returns the decode key. Convert it to the encode key for
        # display. Display on stderr because the output can be redirected.
        print('key = {}'.format(ALEN - key), file=sys.stderr)
        if args.outjson:
            outfile = open(args.outjson, 'w')
            json.dump(fscores, outfile)
    translated = translate(message, key)
    if args.clipboard:
        pyperclip.copy(translated)
    else:
        print(translated)
