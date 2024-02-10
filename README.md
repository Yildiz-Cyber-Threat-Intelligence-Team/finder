# AKINCI Finder Tool

![asciiart](ascii_art.png)


**Parametreler**

    - "arsiv_dosya_yolu": Arama yapılacak arşiv dosyasının veya metin dosyasının yolu.
    - "arama_metni": Arşiv dosyası veya metin dosyası içinde aranacak metin.
    - "dosya_turu": Arşiv dosyasının türü (".zip", ".rar", ".tar.gz", ".tgz"). Metin dosyası için bu parametre kullanılmaz.
    - [cikti_dosyasi]: İsteğe bağlı olarak, sonuçların kaydedileceği dosyanın yolu.

**Kullanım**

    go run main.go "arsiv_dosya_yolu" "arama_metni" "dosya_turu" [cikti_dosyasi]
