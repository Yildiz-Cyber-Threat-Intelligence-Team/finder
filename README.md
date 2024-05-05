# Yıldız Finder Tool

Bu program, belirtilen bir metin dosyasında belirli bir metni arayan basit bir metin dosyası arama programıdır. Kullanıcı, aranacak metni ve aranacak dosyanın yolunu parametre olarak verir. Program, dosyayı açar, her satırı tarar ve aranan metni içeren satırları terminalde gösterir. Eğer aranan metin bulunamazsa, buna dair bir mesaj verir. Arama işlemi tamamlandıktan sonra geçen süreyi de kullanıcıya bildirir.

**Parametreler**

    - "arsiv_dosya_yolu": Arama yapılacak arşiv dosyasının veya metin dosyasının yolu.
    - "arama_metni": Arşiv dosyası veya metin dosyası içinde aranacak metin.
    - "dosya_turu": Arşiv dosyasının türü (".zip", ".rar", ".tar.gz", ".tgz"). Metin dosyası için bu parametre kullanılmaz.
    - [cikti_dosyasi]: İsteğe bağlı olarak, sonuçların kaydedileceği dosyanın yolu.

**Kullanım**

    go run main.go "arsiv_dosya_yolu" "arama_metni" "dosya_turu" [cikti_dosyasi]

**Örnekler**

    Dosya.rar arşiv dosyasında "searchWord" ifadesini aramak için
    ↳ go run main.go "dosya.rar" "searchWord" ".rar"
    Metin.txt dosyasında "arama kelimesi" ifadesini arayıp sonuçları output.txt dosyasına kaydetmek için
    ↳ go run main.go "metin.txt" "arama kelimesi" output.txt
