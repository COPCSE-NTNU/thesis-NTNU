# exerpt from GenData.ipnb

import ftplib
import os

def downloadFTPData(ftp_url, file_name):
    ftp = ftplib.FTP(ftp_url)
    ftp.login()
    ftp.cwd("symboldirectory")
    ftp.retrbinary("RETR " + file_name, open(file_name, "wb").write)
    ftp.quit()


if not os.path.exists("nasdaqlisted.txt"):
    downloadFTPData("ftp.nasdaqtrader.com", "nasdaqlisted.txt")

if not os.path.exists("otherlisted.txt"):
    downloadFTPData("ftp.nasdaqtrader.com", "otherlisted.txt")
