# Terraform провайдер для управления DNS записями на reg.ru

Этот проект содержит Terraform провайдер для управления DNS записями с использованием API reg.ru. Провайдер позволяет создавать, читать и удалять различные типы DNS записей, включая A, AAAA, CNAME, MX и TXT.

## Установка

Для использования этого провайдера вам необходимо установить Terraform версии 0.12 или выше. Вы можете скачать Terraform с [официального сайта](https://www.terraform.io/downloads.html).

## Конфигурация

1. **Создайте файл переменных**:

    Создайте файл `variables.tf` и добавьте следующие переменные:

    ```hcl
    variable "username" {
      description = "Username for the reg.ru API"
      default     = "my_username"
    }

    variable "password" {
      description = "Password for the reg.ru API"
      default     = "my_password"
    }

    variable "cert_file" {
      description = "Path to the client SSL certificate file"
      default     = "./my.crt"
    }

    variable "key_file" {
      description = "Path to the client SSL key file"
      default     = "./my.key"
    }
    ```

2. **Создайте основной конфигурационный файл**:

    Создайте файл `main.tf` с основной конфигурацией для провайдера и ресурсов:

    ```hcl
    terraform {
      required_providers {
        regru = {
          version = "~>0.2.0"
          source  = "letenkov/regru"
        }
      }
    }

    provider "regru" {
      api_username = var.username
      api_password = var.password
      cert_file    = var.cert_file
      key_file     = var.key_file
    }

    resource "regru_dns_record" "example_com" {
      zone   = "example.com"
      name   = "@"
      type   = "A"
      record = "185.199.108.153"
    }

    resource "regru_dns_record" "example_com_ipv6" {
      zone   = "example.com"
      name   = "@"
      type   = "AAAA"
      record = "2606:2800:220:1:248:1893:25c8:1946"
    }

    resource "regru_dns_record" "example_com_mx" {
      zone   = "example.com"
      name   = "@"
      type   = "MX"
      record = "10 mail.example.com"
    }

    resource "regru_dns_record" "example_com_txt" {
      zone   = "example.com"
      name   = "@"
      type   = "TXT"
      record = "v=spf1 include:example.com ~all"
    }
    ```

3. **Инициализация Terraform**:

    В каталоге с конфигурационными файлами выполните команду:

    ```sh
    terraform init
    ```

4. **Планирование конфигурации**:

    Перед применением конфигурации рекомендуется выполнить команду `terraform plan`, чтобы увидеть, какие изменения будут внесены:

    ```sh
    terraform plan
    ```

    Эта команда покажет, какие ресурсы будут созданы, изменены или удалены.

5. **Применение конфигурации**:

    Для создания указанных ресурсов выполните команду:

    ```sh
    terraform apply
    ```

## Разработка и сборка

Для сборки проекта используется `Makefile`. Убедитесь, что у вас установлен Go.

### Шаги по сборке проекта

1. **Клонируйте репозиторий**:

    ```sh
    git clone https://github.com/yourusername/terraform-regru.git
    cd terraform-regru
    ```

2. **Установите зависимости**:

    Выполните команду для установки всех зависимостей:

    ```sh
    make install-deps
    ```

3. **Соберите провайдер**:

    Выполните команду для сборки провайдера:

    ```sh
    make build
    ```

4. **Убедитесь, что провайдер установлен правильно**:

    Проверьте, что собранный провайдер находится в правильной директории:

    ```sh
    ls ~/.terraform.d/plugins/registry.terraform.io/letenkov/regru/0.2.1/$(go env GOOS)_$(go env GOARCH)/
    ```

## Лицензия

Этот проект лицензируется на условиях лицензии Apache 2.0. Подробнее см. в файле [LICENSE](LICENSE).
