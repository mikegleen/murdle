import sys
import string

cipher_06 = 'TVMVIZO XLUUVV WRW MLG SZEV GSV XILDYZI.'
cipher_16 = 'ZTVMG RMP SZW Z NVWRFN-DVRTSG DVZKLM.'
cipher_17 = 'WEEOHRV ASW NO HTE CDEK WSA TEFL HEDNAD.'
cipher_18 = 'NTVIROLSE ETH NGDELE ASW CIPUSISSUO OF HET PENSRO HWO HBTROGU A VHYEA KOBO'
# Silverton the legend was suspicious of the person who brought a heavy book.
cipher = {}
cipher[6] = cipher_06
cipher[16] = cipher_16
cipher[17] = cipher_17
cipher[18] = cipher_18
alph = string.ascii_uppercase
ralph = alph[::-1]
ix = int(sys.argv[1])
lokup = {ralph[i]: alph[i] for i in range(len(ralph))}
# ciphertext = input('Enter ciphertext:').upper()
ciphertext = cipher[ix]
plaintext = ''.join([lokup[c] if c in lokup else c for c in ciphertext])
print(f'{ix=}')
print(ciphertext)
print(plaintext)
