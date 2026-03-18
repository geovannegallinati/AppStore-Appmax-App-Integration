# How to Get Appmax Credentials and Configure Your App

End-to-end guide to create an Appmax App Store application, configure installation URLs, and receive sandbox credentials.

---

## Official References

- Appmax official docs: <https://appmax.readme.io/>
- Create an application: <https://appmax.readme.io/reference/create-an-application>
- App Store portal: <https://appstore.appmax.com.br/>
- CNPJ verification (Receita Federal): <https://solucoes.receita.fazenda.gov.br/Servicos/cnpjreva/cnpjreva_solicitacao.asp>

For daily development, prefer this repository docs first (they are aligned with this codebase and endpoint behavior).

---

## 1. Create Your App Store Account

1. Open <https://appstore.appmax.com.br/> and click `Desenvolva seu aplicativo`.
2. Complete registration at <https://appstore.appmax.com.br/user/create>.
3. Use company data exactly as registered in Receita Federal (CNPJ query):
   - same company email used in CNPJ/company registration
   - exact legal company name
   - exact CNPJ
   - matching company details from Receita Federal
4. Create password, accept terms, and finish registration.
5. Log in at <https://appstore.appmax.com.br/login> with the new email/password.

---

## 2. Create a New Application

After login, use the flow shown by Appmax:

- `Bem-vindo ao ambiente de desenvolvimento da Appmax`
- `Criar um aplicativo`

### Public vs Private App

- `Aplicativo público`: listed in Appmax store; any eligible Appmax user can find and install it.
- `Aplicativo privado`: not listed publicly; only your account and invited users can access/install it.

Choose public when your goal is distribution through the Appmax ecosystem. Choose private for internal/controlled rollout.

---

## 3. Fill Basic App Information

In `Criar aplicativo`, fill:

- `Nome do aplicativo`: clear product name (up to the Appmax limit shown in UI).
- `Email de suporte`: use the same company/CNPJ registration email.
- `Descrição do aplicativo`: short objective description and concrete benefit.
- `Modelo de cobrança`:
  - `Cobrança via Appmax`
  - `Cobrança via Plataforma Externa`

Use the billing model that matches your business flow and contract rules. Define fee/commission terms according to Appmax rules in the official documentation and your commercial agreement.

---

## 4. Upload App Image

In `Imagens do Aplicativo`, upload the icon using Appmax requirements shown in UI:

- square image
- `1200 x 1200` pixels
- `PNG` or `JPG`
- no rounded corners

Use a clean logo with readable contrast; low-quality images can delay review.

---

## 5. Configure Permissions

In `Configurações do Aplicativo`, select webhook/event permissions your integration needs.

Recommended rule:

- enable only events required by your app logic in production
- if you are validating broad webhook compatibility in this project, enable all events used in your test scope

Then save configuration.

---

## 6. Open "My Apps" and Enter Development Mode

1. Go to `Meus aplicativos` (`Consultar aplicativos`).
2. Find the app in status `Desenvolvimento`.
3. Click `Desenvolver`.

If `Desenvolver` does not appear after app creation, send an email to `desenvolvedores@appmax.com.br` asking to enable development fields for your app (it may still be in analysis).

---

## 7. Fill Integration URLs (No Sensitive Domain Exposure)

In the `Desenvolver` modal, configure with your public application domain:

- `Host`: `https://<your-public-domain>/install/start`
- `URL do sistema`: `https://<your-public-domain>/`
- `URL de validação` (Callback URL): `https://<your-public-domain>/integrations/appmax/callback/install`

Example with placeholder:

- `Host`: `https://example-public-domain.com/install/start`
- `URL do sistema`: `https://example-public-domain.com/`
- `URL de validação`: `https://example-public-domain.com/integrations/appmax/callback/install`

Important:

- if callback URL is missing or incorrect, installation will not complete
- ensure HTTPS and a reachable public domain (ngrok/static domain or production domain)

---

## 8. Wait for Appmax Validation and Credential Delivery

Even with all fields configured, development only starts after Appmax validation.

Appmax may validate:

- company identity
- CNPJ status
- ownership/partner records (`quadro societário`)

If approved, Appmax sends sandbox credentials to the company email:

```env
APPMAX_CLIENT_ID=
APPMAX_CLIENT_SECRET=
APPMAX_APP_ID_UUID=
APPMAX_APP_ID_NUMERIC=
```

Only after receiving these values should you finalize `.env` and run full local integration flow.

---

## 9. Next Step in This Repository

After receiving credentials:

1. fill `.env` values
2. follow [Local Development](./local-development.md)
3. use integration docs from this repo:
   - [Architecture Guide](../integration/guide.md)
   - [Endpoints](../integration/endpoints.md)
   - [Frontend Pages](../integration/frontend.md)
