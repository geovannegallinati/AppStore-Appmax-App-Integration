# Appmax Documentation

> Source: https://appmax.readme.io
> Scraped: 2026-03-13
> Pages: 34

---

<!-- url: https://appmax.readme.io/reference/conceitos-de-neg%C3%B3cio -->
## 1. Conceitos de Negócio

Seja bem-vindo à documentação da nossa API.

Esta documentação reúne tudo o que você precisa para integrar com nossa API de forma completa e eficiente, potencializando o seu negócio.

**O que é a Appmax?**

Gateway, antifraude e adquirência em um só lugar para entregar alta performance com simplicidade. A Appmax oferece pagamentos rápidos, seguros e eficientes para impulsionar o seu negócio digital. De e-commerces a SaaS e infoprodutores, milhares de empresas confiam em nossa plataforma para vender mais, com menos esforço e máxima confiabilidade.

Se ainda não conheça a Appmax, convidamos você a visitar nosso site oficial e nossa loja de aplicativos para conhecer melhor nossos produtos e serviços.

---

<!-- url: https://appmax.readme.io/reference/introducao -->
## 1.1. Introdução

Entenda como criar e integrar aplicativos na Appstore da Appmax para expandir recursos, personalizar lojas e conectar sistemas externos.

**O que é a API/Loja de Aplicativos da Appmax?**

A API da Appmax permite criar e gerenciar clientes, pedidos, pagamentos e outros recursos essenciais para o funcionamento de uma loja.
 Ela é pensada para que desenvolvedores integrem de forma eficiente, garantindo segurança e flexibilidade através da criação de um aplicativo. 

A Loja de Aplicativos da Appmax permite que:

Terceiros criem aplicativos para oferecer serviços aos comerciantes.

Comerciantes adicionem funcionalidades extras às suas próprias lojas.

**Benefícios da Appstore:**

- Personalização da experiência da loja.
- Integração com serviços externos.
- Adição de recursos diretamente no painel de administração.
**Tipos de Aplicativos:**

**Aplicativo Privado:** ideal para desenvolvedores que desejam integrar e processar vendas exclusivamente em seus próprios ambientes. Esse tipo de app não fica disponível publicamente na base da Appmax, podendo ser acessado apenas por meio de um link compartilhado pelo próprio desenvolvedor. Essa funcionalidade permite que parceiros criem soluções personalizadas sem expô-las ao público, garantindo uma experiência mais segura, exclusiva e otimizada, tanto para desenvolvedores quanto para lojistas.

**Aplicativo Público:** é um aplicativo que ficará disponível para toda a base de clientes da Appmax. Com ele, desenvolvedores podem oferecer suas soluções de forma ampla, permitindo que qualquer lojista utilize o app diretamente através da plataforma.

**Segurança:**

- O acesso à API é restrito por par de chaves.
- As chaves são emitidas na instalação do aplicativo e devem ser armazenadas de forma segura.
- O uso de tokens temporários garante que credenciais não fiquem expostas em requisições.

---

<!-- url: https://appmax.readme.io/reference/requisitos-tecnicos -->
## 1.1.1. Requisitos Técnicos

## Obrigatórios:

### Autenticação:

É o processo de verificar se quem está tentando acessar a API tem permissão para isso. É uma forma de garantir que apenas usuários ou sistemas autorizados consigam usar os recursos da API. A autenticação com a API da Appmax deve seguir o fluxo de autenticação e autorização descrito na documentação oficial. 

### Segurança Appmax:

 Devido as rigorosas diretrizes do PCI DSS (Padrão de Segurança de Dados da Indústria de Cartões de Pagamento), é essencial implementar o Appmax.js, uma biblioteca JavaScript da Appmax que deve ser incluída no front-end do site sem alterar o visual da loja. Sua principal função é proteger dados sensíveis do cartão, garantindo que essas informações não passem pelos servidores externos. 

### Criação de Cliente:

Para que seja possível criar pedidos na Appmax, primeiro é necessário registrar os dados do comprador (cliente final) no sistema. Esse cadastro gera um identificador único (customer_id), que será utilizado posteriormente na criação dos pedidos, pagamentos e estornos. 

**Atenção:** O IP do cliente final deve ser enviado durante esse processo.

### Criação de Pedido:

Criar um pedido na API significa registrar a compra vinculada a um cliente. Esse processo gera um order_id, que será usado para acompanhar o status do pagamento, realizar estornos e outras operações relacionadas. 

### Pagamentos

- **Cartão de Crédito:** é uma forma de comprar um produto ou serviço usando limite de crédito fornecido pelo emissor do cartão, com a possibilidade de pagar o valor integral ou parcelado depois. Existem duas formas de processar o pagamento via cartão: -**Via Appmax JS (uso de CDN):** é uma forma de capturar os dados do cartão diretamente no navegador do cliente, usando um JavaScript fornecido pela Appmax, para cumprir as regras de segurança PCI DSS (Payment Card Industry Data Security Standard) sem que os dados sensíveis passem pelo seu servidor. - **Via Tokenização (API**):** é um processo usado para substituir os dados reais de um cartão de crédito (como número, CVV e validade) por um token único e seguro, para que transações possam ser feitas sem expor as informações originais.
- **Pix:** é o sistema de pagamento instantâneo criado e gerenciado pelo Banco Central do Brasil. 
- **Boleto:** é basicamente um documento de cobrança emitido por um banco, que contém: Valor a pagar, Data de vencimento, Código de barras e/ou linha digitável, Informações do beneficiário (quem vai receber o pagamento) e Dados do pagador.
- **ApplePay:** é o sistema de pagamentos da Apple que permite fazer compras usando iPhone, Apple Watch, iPad ou Mac, sem precisar digitar dados do cartão toda vez. Funciona com cartões de débito e crédito. 
### Cálculo de Parcelas:

Rota usada para retornar os valores de cada parcela, com base no total do pedido e nas taxas de parcelamento configuradas na Appmax. Evita configurações manuais externas, mantendo valores alinhados às regras da Appmax, suportando de 1 a 12 parcelas, configuráveis diretamente na Appmax. Possui duas modalidades: - **PP (Simples por Parcela):** juros aplicados diretamente em cada parcela. - **AM (Financiamento): **juros aplicados mensalmente sobre o saldo devedor total.

### Estorno de Pagamentos:

 É a devolução de um valor pago ao comprador, feita pelo vendedor ou pela operadora de pagamento, usando o mesmo meio de pagamento utilizado na compra. Existem duas opções para a realização de estornos:

- **Via API:**O merchant envia uma requisição solicitando o estorno diretamente de sua plataforma. Recomendado para **cartão de crédito** e **Pix**, pois no **boleto **é necessário informar os dados bancários do cliente na Appmax para que o valor seja ressarcido.
- **Via Painel Appmax:** O merchant solicita o estorno manualmente no painel da Appmax, será enviado o evento de refund via webhook para atualização do status do pedido na sua plataforma . Recomendado para cartão de credito, pix e boleto. 
### Código de Rastreio:

É um identificador único fornecido pela transportadora ou serviço de entrega, que permite acompanhar a confirmação da entrega do produto encomenda para o cliente. Para que os saques do merchant sejam aprovados, é necessário atualizar o pedido na Appmax com o código de rastreamento da entrega.

### Webhooks:

É um mecanismo de comunicação entre sistemas que permite que um aplicativo envie informações em tempo real para outro quando um evento específico acontece, sem que o outro sistema precise ficar consultando constantemente. A Appmax envia notificações de eventos via Webhook. Todos os eventos são importantes, mas os essenciais para o funcionamento da integração estão descritos em:

## Outras Funcionalidades Opcionais:

### Recorrência:

A recorrência é um modelo de pagamento em que o cliente é cobrado periodicamente por um produto ou serviço, de forma automática, sem precisar autorizar cada pagamento manualmente. _Em fase de beta tester, solicite caso tenha interesse em participar._

### Link de pagamento por API:

É uma URL gerada pela Appmax que permite que um cliente pague um produto ou serviço de forma rápida e segura, sem precisar de um carrinho de compras completo. _Em fase de beta tester, solicite caso tenha interesse em participar._

### Upsell

É uma rota é uma estratégia de vendas usada para incentivar o cliente a comprar um produto complementar do item que está adquirindo ou considerando comprar.

### Recuperação de Vendas com IA

A Appmax oferece uma solução inovadora para recuperar carrinhos de compra abandonados, utilizando inteligência artificial para aumentar a conversão da loja do Merchant. Com nosso endpoint dedicado, suas integração pode enviar automaticamente os carrinhos de abandono, permitindo ações personalizadas de recuperação de vendas. Saiba mais _Em fase de beta tester, solicite caso tenha interesse em participar._

#### Cenários contemplados pela recuperação de vendas:

- Abandono no pagamento com cartão de crédito;
- Abandono no pagamento via Pix;
- Pedidos negados no cartão de crédito;

---

<!-- url: https://appmax.readme.io/reference/diferencas-entre-os-ambientes-da-appmax-sandbox-e-producao -->
## 1.1.2. Diferença entre os ambientes de Sandbox e Produção

Entenda a diferença entre os ambientes disponíveis no desenvolvimento do seu aplicativo na Appmax.

## Diferença entre os ambientes Sandbox e Produção

A Appmax oferece dois ambientes distintos para integração:

####  Sandbox (Ambiente de Testes)

- Para testar a integração antes de ir para produção e de ser aprovado na homologação.
- A autenticação e a API utilizam subdomínios diferentes: 
  - Autenticação Sandbox: 
```
https://auth.sandboxappmax.com.br
```

  - API Sandbox: 
```
https://api.sandboxappmax.com.br
```

- O redirecionamento para autorização ocorre na **Appmax BC Sandbox**: 
```
https://breakingcode.sandboxappmax.com.br/appstore/integration/HASH
```

####  Produção

- Para transacionar com clientes reais.
- URLs padrão de autenticação e API: 
  - Autenticação Produção: 
```
https://auth.appmax.com.br
```

  - API Produção: 
```
https://api.appmax.com.br
```

- O redirecionamento para autorização ocorre no **Admin da Appmax**: 
```
https://admin.appmax.com.br/appstore/integration/HASH
```

#### Resumo:

- No sandbox, as URLs começam com **sandboxappmax** e o redirecionamento é para **BC Sandbox**.
- Em produção, as URLs são padrão **appmax** e o redirecionamento é para o **Admin da Appmax**.

---

<!-- url: https://appmax.readme.io/reference/status-de-pedidos -->
## 1.1.3. Status de pedidos

Ao longo do processo de um pedido, ele passa por diferentes status. Abaixo, há uma listagem que descreve cada um desses status, os quais representam fases específicas do processo do pedido. Além disso apresentaremos os status que você deverá usar nas suas integrações, mas ao entrar no painel, eles podem ter outra nomenclatura, visando uma experiência melhor do lojista.

---

**pendente**

Status no painel: Pagamento pendente

Todos os pedidos que ainda não estão pagos, um pedido de cartão de crédito que ainda não foi autorizado, um pedido de pix que ainda não foi pago ou um pedido de boleto que não compensou, isso inclui os de boletos vencidos que nunca foram pagos.

---

**aprovado**

Status no painel: Pagamento aprovado

Indica que está tudo certo com o pagamento do pedido, a partir desse status os valores já são disponibilizados na conta do lojista.

---

**autorizado**

Status no painel: Análise antifraude

Um pedido de cartão de crédito passa do status pendente para esse status ao ser autorizado pelo banco emissor e possuir saldo, a partir desse momento é iniciado o processo de análise antifraude e se tudo estiver correto, o pedido irá para aprovado.

---

**cancelado**

Status no painel: Não autorizado

Pedidos de PIX que tiveram o tempo de validade do QRCode expirado, pedidos de cartão que não foram autorizados pelo banco emissor, ou não possuem saldo ficam nesse status.

Há uma exceção também para pedidos criados internamente pelo painel, que nunca receberam uma transação de pagamento também se encontram com esse status.

---

**estornado**

Status no painel: Estornado

Um pedido que recebeu uma solicitação de reembolso aprovada por parte do cliente ou por parte do lojista recebe esse status, no momento que um pedido passa por esse status é realizado o débito do valor na conta do lojista.

---

**recusado_por_risco**

Status no painel: Recusado por risco

Para transações que são consideradas de alto risco e não passam na análise de antifraude, o pedido será estornado e depois receberá esse status.

---

**integrado**

Status no painel: Pagamento aprovado

Esse é o status final para um pedido aprovado, após passar pelas validações das integrações, ele está pronto para ser enviado ao comprador final (em casos de produtos físicos).

---

**pendente_integracao**

Status no painel: Pagamento aprovado

Esse status indica que o pedido está pago, mas há alguma pendência com a integração, isso pode acontecer por alguma informação incorreta no cadastro ou no produto.

---

**pendente_integracao_em_analise**

Status no painel: Pagamento aprovado

Quando um pedido foi aprovado, mas antes de passar para integrado ele recebe uma solicitação de estorno, esse status é atribuído, a partir desse momento uma análise manual será executada, pois caso a solicitação de estorno seja aprovada, os produtos não devem ser enviados, esse status é incomum mas pode acontecer.

---

**chargeback_em_tratativa**

Status no painel: Chargeback

Quando recebemos uma sinalização do banco que um chargeback ocorreu, o pedido é estornado e depois recebe esse status. Saiba mais neste artigo.

---

**chargeback_em_disputa**

Status no painel: Chargeback

Quando a Appmax inicia uma disputa para recuperar um chageback, o pedido sai do status de tratativa e recebe esse status.

---

**chargeback_perdido**

Status no painel: Chargeback

Quando não há mais formas de recuperar um chargeback, o pedido recebe esse status.

---

**chargeback_vencido**

Status no painel: Chargeback recuperado

Quando vencemos uma disputa de chargeback, o pedido receberá esse status, a partir desse momento, o lojista recebe o crédito em seu saldo.

---

<!-- url: https://appmax.readme.io/reference/sla-e-orienta%C3%A7%C3%B5es-para-homologa%C3%A7%C3%A3o -->
## 1.1.4. SLA e orientações para homologação

O processo de homologação possui uma série de etapas e cenários essenciais para validar a funcionalidade e a conformidade dos aplicativos na Loja de Aplicativos. Abaixo está um resumo organizado e ajustado para facilitar a compreensão e a aplicação.

**Processo de Homologação**

**Objetivo**

Validar os dados do aplicativo e realizar testes ponta a ponta para garantir que ele atenda aos padrões exigidos.

---

**Início do Processo**

É iniciado até **3 dias úteis** após a integração sinalizar a finalização da implementação.
 O processo pode durar até **7 dias úteis**.

---

**Regra para Adequações**

A cada nova solicitação de adequação, o aplicativo retorna para o final da fila de homologação e com o mesmo prazo citado acima

---

**Requisitos para Homologação**

- **Sinalizar no canal oficial** com a equipe de integração.
- **Disponibilizar uma conta e uma loja de teste**, para poder simular o processo igual a um merchant.
- **Enviar os acessos necessários** para o e-mail: [email protected].
- Disponibilizar uma documentação com o passo a passo da instalação do aplicativos.
Segue abaixo a lista de cenários essenciais que serão avaliados durante a homologação do aplicativo na Appmax. Porém, outros testes podem ser realizados, mas recomendamos que os listados abaixo sejam verificados previamente para aumentar a assertividade no processo:

---

**Cenários de Testes:**

**Validações Gerais**

- Validação do logo do aplicativo.
- Validação da descrição do aplicativo.
- Validação do e-mail de suporte.
- Instalação do aplicativo.
**Testes de Compra e estorno (Pessoa Física - PF)**

- Cenário de compra no cartão com juros com CPF.
- Cenário de compra no cartão sem juros com CPF.
- Cenário de compra no Pix com CPF.
- Cenário de compra no boleto com CPF.
- Estorno total no cartão com juros com CPF.
- Estorno total no cartão sem juros com CPF.
- Estorno total no Pix com CPF.
- Estorno total no boleto com CPF.
- Estorno parcial no cartão com juros com CPF (via Appmax).
- Estorno parcial no cartão sem juros com CPF (via Appmax).
**Testes de Compras e Estornos (Pessoa Jurídica - PJ)**

- Compra no Pix com CNPJ.
- Compra no boleto com CNPJ.
- Compra no cartão com juros com CNPJ.
- Compra no cartão sem juros com CNPJ.
- Estorno total no Pix com CNPJ.
- Estorno total no boleto com CNPJ.
- Estorno total no cartão sem juros com CNPJ.
- Estorno total no cartão com juros com CNPJ.
- Estorno parcial no cartão com juros com CNPJ.
- Estorno parcial no cartão sem juros com CNPJ.
- Estorno parcial no Pix com CNPJ.
- Estorno parcial no boleto com CNPJ.
**Outros Cenários**

- Avaliar como ficará o softer descriptor nas compras de cartão.
- Testar integração do código de rastreio.
- Avaliar o IP registrado nos pedidos.
- Testar compras com cupom de desconto.
- Testar compras com frete e juros.
- Testar compras com frete + cupom de desconto.
- Testar compras com mais de um produto diferente no carrinho.
- Testar compras com múltiplas unidades do mesmo produto no carrinho.
- Testar a atualização de status dos pedidos.
Obs.: Se o seu app não for voltado para o processamento de pagamentos pela Appmax, os cenários acima podem ser desconsiderados e a análise será realizada com base em outros critérios.

---

<!-- url: https://appmax.readme.io/reference/autentica%C3%A7%C3%A3o-e-autoriza%C3%A7%C3%A3o-na-api -->
## 1.1.5. Autenticação e Autorização na API

**Motivo para Não Utilização de Refresh Tokens**

Nossa API adota um modelo de autenticação e autorização que não utiliza refresh tokens. Esta decisão baseia-se na arquitetura específica da nossa integração, que é projetada para comunicação entre servidores (server-to-server), ao invés de uma aplicação frontend (como uma Single Page Application - SPA) se comunicando com um backend. Abaixo estão os principais pontos que justificam essa abordagem:

1. **Natureza da Integração Server-to-Server:** 
  - Em nossa arquitetura, a comunicação ocorre diretamente entre servidores, o que significa que os tokens de acesso são gerenciados em ambientes controlados e seguros. Isso reduz a necessidade de mecanismos adicionais para a renovação de tokens, como os refresh tokens.
2. **Segurança e Simplicidade:** 
  - A utilização de tokens de acesso de curta duração oferece uma camada robusta de segurança ao limitar a janela de tempo em que um token pode ser utilizado. Em um ambiente server-to-server, onde as credenciais e os tokens podem ser armazenados de maneira segura, esta abordagem simplifica o gerenciamento de tokens sem comprometer a segurança.
3. **Redução da Complexidade:** 
  - Implementar e gerenciar refresh tokens adiciona uma camada extra de complexidade. Isso inclui o armazenamento seguro de refresh tokens, a rotação de tokens, e a lógica para a renovação de tokens de acesso. Dado que nossos servidores são ambientes confiáveis e seguros, optamos por uma abordagem mais direta e eficaz.
4. **Conformidade com Melhores Práticas de Server-to-Server:** 
  - Em integrações server-to-server, é comum a utilização de tokens de acesso curtos com autenticação baseada em chaves ou certificados. Este método é amplamente aceito como uma prática segura e eficiente, que atende bem às necessidades de comunicação entre servidores.
**Como Funciona:**

- **Autenticação Inicial:** Os servidores autenticam-se utilizando chaves de API ou certificados digitais.
- **Tokens de Acesso de Curta Duração:** Após a autenticação, um token de acesso de curta duração é emitido e deve ser usado para todas as requisições subsequentes.
- **Renovação de Tokens: **Quando um token de acesso expira, o servidor pode obter um novo token através do processo de autenticação inicial, mantendo a segurança e a simplicidade.
**Vantagens:**

- **Segurança: **Limita o período de uso dos tokens de acesso, reduzindo o risco de utilização indevida em caso de comprometimento.
- **Simplicidade:** Elimina a necessidade de armazenamento e gerenciamento de refresh tokens.
- **Eficiência: **Adequado para ambientes controlados e seguros, onde a comunicação ocorre diretamente entre servidores.
**Nossa abordagem assegura que a autenticação e autorização na nossa API sejam realizadas de maneira segura e eficiente, adequada à natureza de integrações server-to-server.**

---

<!-- url: https://appmax.readme.io/reference/simulador-de-cart%C3%A3o-de-cr%C3%A9dito -->
## 1.1.6. Simulador de cartão de crédito

Para simular transações de cartão de crédito definimos algumas regras que devem ser utilizadas. Cada uma dessas regras implica em uma resposta da API, e dessa forma você pode testar transações de cartão de crédito.

Para testar cada cenário abaixo, é preciso enviar o respectivo número de cartão, com uma data de expiração futura.

---

**Os números dos cartões e seus respectivos cenários são esses:**

| Número do cartão | Cenário |
| --- | --- |
| 4000000000000010 | Cartão de sucesso. Qualquer operação com esse cartão é realizada com sucesso. |
| 4000000000000028 | Cartão de falha. Qualquer transação retorna como "não autorizada". |
| Qualquer outro cartão | Qualquer transação retorna como "não autorizada". |

---

<!-- url: https://appmax.readme.io/reference/webhooks -->
## 1.2. Webhooks

### Como funciona

Quando um dos eventos selecionados ocorrer para o merchant, enviaremos uma requisição para o host que você forneceu durante a criação do aplicativo. Cada requisição terá um payload específico, variando de acordo com o evento que foi desencadeado.

| Evento (Descrição) | event | event_type |
| --- | --- | --- |
| Cliente criado | customer_created | customer |
| Cliente Interessado | customer_interested | customer |
| Cliente Contatado | customer_contacted | customer |
| Pedido Autorizado | order_authorized | order |
| Pedido Aprovado | order_approved | order |
| Boleto Criado | order_billet_created | order |
| Pedido Pago | order_paid | order |
| Pedido Pendente Integração | order_pending_integration | order |
| Pedido Estornado | order_refund | order |
| Upsell Pago | order_up_sold | order |
| Pix Gerado | order_pix_created | order |
| Pix Pago | order_paid_by_pix | order |
| Pix Expirado | order_pix_expired | order |
| Pedido Integrado | order_integrated | order |
| Pedido com Boleto Vencido | order_billet_overdue | order |
| Pedido Autorizado com Atraso | order_authorized_with_delay | order |
| Pedido em Chargeback em Tratamento | order_chargeback_in_treatment | order |
| Assinatura Cancelada | subscription_cancelation | subscription |
| Assinatura Atrasada | subscription_delayed | subscription |
| Pagamento Autorizado com Atraso | payment_authorized_with_delay | payment |
| Pagamento Não Autorizado | payment_not_authorized | payment |

---

<!-- url: https://appmax.readme.io/reference/appmax-js -->
## 1.3. Appmax JS

Entenda como implementar o script de segurança da Appmax

O `appmax.js` é uma biblioteca JavaScript desenvolvida pela Appmax para integração fácil e segura em páginas de checkout. Ao incluir o `appmax.js` em sua página, você não precisa se preocupar com alterações no visual da loja, já que o script atua de forma discreta e eficiente.

## Como funciona

Devido às rigorosas diretrizes do **Padrão de Segurança de Dados da Indústria de Cartões de Pagamento (PCI DSS)**, é crucial proteger dados sensíveis, como informações de cartões de pagamento. Além disso realizamos a coleta de IP para a segurança do fluxo de pagamento. 

O `appmax.js` foi projetado para atender a essas preocupações de segurança. Ao integrar o `appmax.js`, você evita que dados sensíveis do cartão de crédito passem diretamente pelos seus servidores.

O uso desse script é **obrigatório**, logo para integrar com a Appmax, você precisará incluir o script em seu checkout, caso você possua certificação PCI, a funcionalidade de tokenização é opcional, mas a funcionalidade de coleta de IP é obrigatória, através dela conseguimos garantir a segurança das transações e análises em nosso fluxo de pagamento.

## Como usar

1. Incluir o Script da CDN: 
```
<script src="https://scripts.appmax.com.br/appmax.min.js"></script>
```

2. Inicializar o **AppmaxScripts**
 Após carregar o script, você precisa inicializar o **AppmaxScripts** com uma função de sucesso e uma função de erro.
## Funcionalidades disponíveis

1. Coleta do IP do Cliente
2. Tokenização de Pagamento

---

<!-- url: https://appmax.readme.io/reference/criar-app -->
## 2. Criar aplicativo

Entenda como criar um aplicativo na Appmax

Na Appmax, durante a criação de um aplicativo, é possível optar entre dois tipos: **Público** ou **Privado**.

**Aplicativo Privado:** ideal para desenvolvedores que desejam integrar e processar vendas exclusivamente em seus próprios ambientes. Esse tipo de app não fica disponível publicamente na base da Appmax, podendo ser acessado apenas por meio de um link compartilhado pelo próprio desenvolvedor. Essa funcionalidade permite que parceiros criem soluções personalizadas sem expô-las ao público, garantindo uma experiência mais segura, exclusiva e otimizada, tanto para desenvolvedores quanto para lojistas.

**Aplicativo Publico:** é um aplicativo que ficará disponível para toda a base de clientes da Appmax. Com ele, desenvolvedores podem oferecer suas soluções de forma ampla, permitindo que qualquer lojista utilize o app diretamente através da plataforma.

## Criação do aplicativo

Durante o processo de criação, é necessário escolher se o seu aplicativo será **Público** ou **Privado**.

> [Image: This image shows a UI screen for selecting application visibility settings, presenting two radio button options: "Aplicativo público" (Public Application) which allows anyone to install and use the app, and "Aplicativo privado" (Private Application) which keeps the app unlisted from the public store with access only to invited users.]

 Abaixo, veja as principais diferenças entre eles:

| Característica | Aplicativo Público | Aplicativo Privado |
| --- | --- | --- |
| Exibição na AppStore | Visível para todos os usuários | Não aparece na listagem, acesso somente pelo link exclusivo |
| Forma de Acesso | Disponível na página oficial da Loja de Aplicativos da Appmax | Acesso apenas por link de compartilhamento |
| Instalação | Instalação direta pela Loja de aplicativos da Appmax | Instalado através de link, visível para o desenvolvedor após criação do aplicativo possibilitando o compartilhamento |

Para criar um aplicativo, será necessário informar algumas informações sobre seu aplicativo, que será dividido em três etapas:

1. **Sobre o aplicativo**
 Neste espaço, você deverá adicionar algumas informações detalhadas sobre o aplicativo. 
  - **Nome do aplicativo:** O nome que será exibido para os merchants e como utilizaram para pesquisar seu app. Possui limite de até 30 caracteres.
  - **E-mail de suporte:** E-mail que os usuários poderão te acionar em caso de dúvidas.
  - **Descrição do aplicativo: **Deverá descrever sobre objetivo, benefícios e serviços oferecidos pelo aplicativo para que os merchants consigam identificar a funcionalidade de seu aplicativo. Esse texto possui um limite de até 100 caracteres e será exibido junto com o nome.
  - **Modelo de cobrança:** Como deseja que seu aplicativo seja cobrado dos merchants, sendo dividido em duas opções:
  - **Cobrança via Plataforma Externa:** A cobrança será realizada diretamente em sua plataforma, sem envolvimento da Appmax.
  - **Cobrança via Appmax:** Nessa opção, você deverá definir um valor fixo que será cobrado do merchant mensalmente na Appmax, descontando o valor diretamente do saldo parceiro. Nesse formato, o time entrará em contato com você para solicitar seus dados bancários para que seja realizado o repasse de seus valores.
---

1. **Imagens do Aplicativo**
 Faça o upload da imagem que servirá como avatar para o seu aplicativo. É importante observar que a imagem deve atender aos requisitos de tamanho e tipo de arquivos exigidos.
 **A imagem deve ser um quadrado no tamanho 1200px x 1200px, PNG ou JPG, sem cantos arredondados.**
---

1. **Configurações do Aplicativo**
 Nesta etapa, você terá a oportunidade de escolher quais eventos de webhook serão enviados para o seu aplicativo.
---

**Após a configuração do aplicativo**
 Você será direcionado para uma página de sucesso e deverá clicar na opção “Consultar Aplicativo” para concluir o desenvolvimento de seu aplicativo.
 Será exibida uma lista, e você deverá clicar na opção “Desenvolver“, o que abrirá uma modal.

- **Id do aplicativo:** A Chave de API que você deverá utilizar para se comunicar com o sistema da Appmax. Essa chave serve como uma credencial de segurança, permitindo que o seu aplicativo autentique e acesse os recursos necessários na plataforma da Appmax. Certifique-se de proteger essa chave e utilizá-la de forma segura para garantir a integridade e segurança das suas interações com o sistema.
- **Host:** é necessário fornecer a URL para onde enviaremos as requisições dos eventos escolhidos. Esta URL servirá como o endpoint para os webhooks do seu aplicativo, permitindo que receba notificações em tempo real sobre os eventos selecionados. Certifique-se de fornecer uma URL válida e segura para garantir a correta entrega das requisições.
- **URL do sistema:** Adicione a URL do seu sistema. Esta é a URL que será disponibilizada para o merchant conseguir acessar seu sistema.
- **URL de validação:** Adicione uma URL para a instalação do aplicativo. Será utilizada para validar que o aplicativo foi instalado com sucesso pelo merchant.
**Observação:** Neste modal, você tem duas opções disponíveis. Ao selecionar “Enviar para análise“, o aplicativo será direcionado para análise por nossa equipe. Se optar por “Testar“, receberemos uma notificação e entraremos em contato para fornecer os acessos necessários em nossa área de homologação para realizar os testes.

Fique atento ao seu e-mail cadastrado. Qualquer comunicação referente ao status da análise do seu aplicativo será enviada para esse endereço de e-mail. Certifique-se de verificar regularmente sua caixa de entrada, incluindo a pasta de spam ou lixo eletrônico, para garantir que você não perca nenhuma atualização importante sobre o processo de análise do seu aplicativo na Appmax.

---

<!-- url: https://appmax.readme.io/reference/sobre-as-credenciais -->
## 2.2. Sobre as credenciais

Entenda como funciona as credenciais do aplicativo e como obter as credenciais da API da Appmax.

### Sobre as credenciais

A API da Appmax utiliza **dois pares de credenciais** (`client_id` e `client_secret`), cada um com uma finalidade diferente:

1. **Credenciais do Aplicativo (App credentials)** 
  1. Essas credenciais (`client_id` e `client_secret`) são utilizadas apenas para a instalação do aplicativo e geração das credenciais do **merchant**.
  2. O escopo dessas credenciais é **limitado à autenticação e autorização da instalação**, não permitindo a realização de transações na API.
  3. Elas devem ser utilizadas em conjunto com o `app_id` para autorizar a instalação.
2. **Credenciais do Merchant (Merchant credentials)** 
  1. Para **cada instalação do aplicativo**, será gerado **um novo par de credenciais** (`client_id` e `client_secret`) exclusivo para o **merchant**.
  2. Essas credenciais têm escopo completo na API com base nas configurações do seu aplicativo, permitindo a criação de **clientes, pedidos, pagamentos e outras operações transacionais**, se for o caso.
  3. Após a instalação e autorização do aplicativo, esse par de credenciais será utilizado para todas as requisições à API.
> [Image: # Sequence Diagram: AppMax API Integration Flow

This sequence diagram illustrates a 5-step authentication and authorization flow between an Integration client and three AppMax services (auth.appmax.com.br, admin.appmax.com.br, and api.appmax.com.br), showing how initial credentials are exchanged for API tokens that enable payment endpoint consumption.]

---

<!-- url: https://appmax.readme.io/reference/instalacao-do-aplicativo -->
## 2.3. Instalação do Aplicativo

Entenda como obter o token do aplicativo e realizar a autorização da instalação.

## 1. Obter o token do aplicativo

Antes de qualquer requisição, é necessário obter o **token de acesso** do seu aplicativo.

#### Endpoint:

`POST` https://auth.appmax.com.br/oauth2/token

##### Requisição (cURL)

```
curl --location 'https://auth.appmax.com.br/oauth2/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=client_credentials' \
--data-urlencode 'client_id=CLIENT_ID' \
--data-urlencode 'client_secret=CLIENT_SECRET'
```

##### Exemplo de Resposta

```
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6Ikp...",
    "token_type": "Bearer",
    "expires_in": 3600
}
```

**Importante:** O `access_token` obtido aqui será utilizado nas próximas requisições como **Bearer Token**.

## 2. Autorizar a instalação do aplicativo

Agora, com o token de acesso **do aplicativo**, precisamos gerar um `hash` de autorização para fazer o redirect e criar as credenciais do **merchant**.

- **app_id:** Chave de API utilizada para que seu aplicativo se comunique com o sistema da Appmax. Ela funciona como uma credencial de segurança, permitindo autenticação e acesso aos recursos necessários na plataforma.
- **external_key:** Chave fornecida pelo sistema da plataforma parceira, usada para identificar com precisão a origem da instalação. Normalmente corresponde a algum identificador já existente no ambiente do cliente/lojista (como store_id, merchant_id, company_key etc.). Essa informação não é controlada pela Appmax, apenas recebemos e repassamos para que a plataforma utilize e correlacione com o sistema dela.
- **url_callback:** URL para onde o usuário será redirecionado após a geração e validação do hash, permitindo que a plataforma finalize o fluxo de criação das credenciais do merchant.
##### Requisição para gerar o hash (token):

```
curl --location 'https://api.appmax.com.br/app/authorize' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer SEU_TOKEN' \
--data '{
    "app_id": "APP_ID",
    "external_key": "EXTERNAL_KEY",
    "url_callback": "URL_CALLBACK"
}'
```

##### Exemplo de resposta com o hash:

```
{
    "data": {
        "token": "12083w36219d223f33ecf48f2a7f5ccf143b0bc554"
    }
}
```

> 
### 
O token recebido é o hash que deve ser utilizado para redirecionar o usuário ao sistema da Appmax para autorizar a instalação.

#####  URL de Redirecionamento

> 
### 
Substitua o "HASH" no final da URL pelo hash gerado anteriormente

**Exemplo:** https://breakingcode.sandboxappmax.com.br/appstore/integration/12083w36219d223f33ecf48f2a7f5ccf143b0bc554

- **Sandbox**: https://breakingcode.sandboxappmax.com.br/appstore/integration/HASH
- **Produção**: https://admin.appmax.com.br/appstore/integration/HASH

---

<!-- url: https://appmax.readme.io/reference/criar-credenciais-do-merchant -->
## 2.4. Criar as credenciais do Merchant

Entenda o fluxo de criação das credenciais do Merchant e como implementar o Health Check de instalação.

## 1. Criar credenciais do Merchant

Após o usuário autorizar a instalação, o próximo passo é gerar as credenciais do **merchant**, que serão usadas para transacionar na API (criar pedidos, clientes, pagamentos, etc).

#### Endpoint:

`POST` https://api.appmax.com.br/app/client/generate

##### Requisição (cURL

```
curl --location 'https://api.appmax.com.br/app/client/generate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer SEU_TOKEN' \
--data '{
    "token": "12083w36219d223f33ecf48f2a7f5ccf143b0bc554"
}'
```

##### Exemplo de Resposta

```
{
    "data": {
        "client": {
            "client_id": "MERCHANT_CLIENT_ID",
            "client_secret": "MERCHANT_CLIENT_SECRET"
        }
    }
}
```

Essas credenciais são exclusivas do **merchant** e devem ser usadas para realizar transações na API.

#### Health check no aplicativo para concluir a instalação:

Quando houver a instalação do aplicativo por parte do **merchant** na Appmax, será enviado uma requisição `POST` para a URL especificada informada no campo **URL de validação**, localizado no modal durante a etapa de criação do aplicativo. Esta requisição conterá o **app_id** (ID do aplicativo), **client_id** e **client_secret** que são as chaves de acesso para API da Appmax, e a **external_key** que é uma chave que os usuários que instalam seu aplicativo podem fornecer para uma identificação mais precisa:

```
{
  "app_id": "APP_ID",
  "client_id": "MERCHANT_CLIENT_ID",
  "client_secret": "MERCHANT_CLIENT_SECRET",
  "external_key": "EXTERNAL_KEY"
}
```

Esperamos como resposta um retorno com **HTTP Code 200**, contendo a resposta que inclui o `external_id`. Este `external_id` será um **ID único** para cada instalação armazenado em seu banco de dados, facilitando a identificação do vínculo criado em sua base.

```
{
  "external_id": "37bb0791-ee0b-457d-860c-186e32978bcd"
}
```

#### Resumo:

- **Obtenha o token do aplicativo** usando as credenciais do seu app (`client_id` e `client_secret`).
- **Autorize a instalação do aplicativo**, gerando um `hash` que será utilizado no redirecionamento.
- **Após a autorização do usuário**, utilize o `hash` para **gerar as credenciais do merchant**.
- **Utilize essas credenciais (`client_id` e `client_secret`) para realizar transações na API.**

---

<!-- url: https://appmax.readme.io/reference/criar-ou-atualizar-cliente-na-appmax -->
## 3. Como criar ou atualizar um Cliente

Entenda como fazer a criação ou a atualização de um cliente na Appmax

Esta rota permite criar ou atualizar um cliente com ou sem produto vinculado, utilizando como chave de identificação a combinação:

 `first_name` + `last_name` + `email` + `phone` + `ip`.

> 
### 
**Dica:**

Se enviar apenas os campos obrigatórios, o cliente será registrado como **"carrinho abandonado"**, podendo ser atualizado depois pela mesma rota.

> 
### 
Pré-requisito:

Para criar um cliente, você** precisa ter** feito a coleta de IP utilizando o script Appmax JS.

 Caso ainda não tenha feito a coleta de IP, siga antes a documentação: 

Coletar IP do Cliente

### Exemplo de Requisição (cURL)

```
curl --request POST \
     --url https://api.appmax.com.br/v1/customers \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
  "first_name": "Junior",
  "last_name": "Almeida",
  "email": "[email protected]",
  "phone": "51983655100",
  "document_number": "25226493029",
  "address": {
    "postcode": "91520270",
    "street": "Rua Francisco Carneiro da Rocha",
    "number": "582",
    "complement": "Casa",
    "district": "Moinhos de Ventos",
    "city": "Porto Alegre",
    "state": "RS"
  },
  "ip": "127.0.0.1",
  "products": [
    {
      "sku": "9000010",
      "name": "Livro de receitas",
      "quantity": 1,
      "unit_value": 12300,
      "type": "digital"
    }
  ],
  "tracking": {
    "utm_source": "google",
    "utm_campaign": "teste"
  }
}
'
```

### Parâmetros do Body:

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| first_name | string |  | Nome do cliente |
| last_name | string |  | Sobrenome do cliente |
| email | string |  | E-mail válido |
| phone | string |  | Telefone com DDD |
| document_number | string |  | CPF ou CNPJ |
| address | object |  | Endereço do cliente |
| ip | string |  | IP de origem |
| products | array |  | Lista de produtos vinculados |
| tracking | object |  | Dados de origem da visita |
`````````````````` 
## Exemplo de Resposta (JSON)

#### `201` – Cliente criado com sucesso

```
{
  "data": {
    "customer": {
      "id": 1
    }
  }
}
```

#### `422` – Erro de validação dos dados

```
{
  "message": "The given data failed to pass validation.",
  "errors": {
    "message": {
      "first_name": [
        "The first_name field is required."
      ],
      "last_name": [
        "The last_name field is required."
      ],
      "phone": [
        "The phone field must be a string.",
        "The phone field must have a maximum of 11 characters."
      ],
      "email": [
        "The email field is required.",
        "The email field must be a string."
      ],
      "ip": [
        "The ip field is required."
      ]
    }
  }
}
```

> 
### 
Importante:

Guarde o valor de `customer_id` retornado nesta etapa, mesmo que temporariamente.

Ele será necessário para criar o pedido na próxima etapa do fluxo.

** Sem este ID, não será possível criar um pedido pois é necessário vincular o cliente ao pedido.**

---

<!-- url: https://appmax.readme.io/reference/coleta-de-ip-do-customer -->
## 3.1. Coleta de IP do Customer

Como fazer a coleta de IP utilizando o Appmax JS para criar um Cliente na Appmax.

Ao utilizar o atributo **data-appmax-customer** em um formulário, a CDN coleta automaticamente o endereço IP do cliente. Isso elimina a necessidade de lógica adicional para capturar o IP, simplificando o processo de coleta de dados do cliente.

#### Exemplo de Formulário com Coleta de IP:

O IP pode ser adicionado automaticamente no formulário HTML

```
<form id="customer-form" data-appmax-customer>
  <div>
    <label for="first-name">Primeiro Nome:</label>
    <input type="text" id="first-name" name="first-name" required />
  </div>
  <div>
    <label for="last-name">Último Nome:</label>
    <input type="text" id="last-name" name="last-name" required />
  </div>
  <div>
    <label for="email">Email:</label>
    <input type="email" id="email" name="email" required />
  </div>
  <div>
    <label for="phone">Telefone:</label>
    <input type="text" id="phone" name="phone" required />
  </div>
  <div>
    <label for="document-number">Número do Documento:</label>
    <input type="text" id="document-number" name="document-number" required />
  </div>
  <button type="submit">Enviar</button>
</form>
```

Ou também, ao utilizar frameworks, recuperando os dados no método de sucesso e enviando no formulário de criação do customer

```
<script setup>
import { ref, onMounted } from 'vue'

const customer = ref({
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  ip: ''
})

const success = (data) => {
  customer.value.ip = data.ip || 'IP não encontrado.'
}

const error = (error) => {
  console.error('Error:', error)
}

onMounted(() => {
  const script = document.createElement('script')
  script.src = 'https://scripts.appmax.com.br/appmax.min.js'
  script.onload = () => {
    if (window.AppmaxScripts) {
      window.AppmaxScripts.init(success, error)
    } else {
      console.error('AppmaxScripts não carregado.')
    }
  }
})

const submitForm = async () => {
  // Enviar o objeto customer
}
</script>

<template>
  <div>
    <h1>Customer Form</h1>
    <form @submit.prevent="submitForm" data-appmax-customer>
      <div>
        <label for="firstName">First Name:</label>
        <input v-model="customer.first_name" id="firstName" type="text" required />
      </div>
      <div>
        <label for="lastName">Last Name:</label>
        <input v-model="customer.last_name" id="lastName" type="text" required />
      </div>
      <div>
        <label for="email">Email:</label>
        <input v-model="customer.email" id="email" type="email" required />
      </div>
      <div>
        <label for="phone">Phone:</label>
        <input v-model="customer.phone" id="phone" type="text" required />
      </div>

      <button type="submit">Submit</button>
    </form>
  </div>
</template>
```

---

<!-- url: https://appmax.readme.io/reference/como-criar-um-pedido-na-appmax -->
## 4. Como criar um pedido na Appmax

Nesta seção, você aprenderá o fluxo completo para criar um pedido na Appmax — desde o cadastro do cliente até a consulta final do status do pedido.
 Esse processo é dividido em quatro etapas, e cada uma delas possui sua própria subseção com detalhes, exemplos e boas práticas.

###  Visão Geral do Fluxo

1. **Criar ou atualizar um cliente**
 Cadastra um novo cliente ou atualiza um existente para vinculá-lo ao pedido.
2. **Criar o pedido**
 Gera o registro do pedido no sistema, associando-o ao cliente.
3. **Efetuar o pagamento**
 Processa o pagamento usando um dos métodos disponíveis: Apple Pay, cartão de crédito, Pix ou boleto.
4. **Consultar o status e os dados do pedido**
 Recupera informações atualizadas sobre o pedido, incluindo status de pagamento e entrega.

---

<!-- url: https://appmax.readme.io/reference/criar-um-um-pedido-na-appmax -->
## 4.1. Criar um Pedido

Como criar um Pedido na Appmax.

O endpoint `POST /v1/orders` é utilizado para criar um novo pedido na Appmax.

 Um pedido deve estar sempre vinculado a um **cliente previamente criado** no sistema.

> 
### 
Pré-requisito:

Para criar um pedido, você** precisa ter** o `customer_id` do cliente.

 Caso ainda não tenha criado o cliente, siga antes a documentação: 

Criar ou atualizar cliente na Appmax

### Exemplo de Requisição (cURL)

```
curl --request POST \
     --url https://api.appmax.com.br/v1/orders \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
  "customer_id": 29,
  "products_value": 12300,
  "discount_value": 0,
  "shipping_value": 3999,
  "products": [
    {
      "sku": "9000010",
      "name": "Livro de receitas",
      "quantity": 1,
      "unit_value": 12300,
      "type": "digital"
    }
  ]
}
'

```

### Parâmetros do Body:

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| customer_id | number |  | ID do cliente (obtido na etapa de criação do cliente) |
| products_value | number |  | Valor total dos produtos |
| discount_value | number |  | Valor do desconto |
| shipping_value | number |  | Valor do frete |
| products | array |  | Lista de produtos do pedido |
`````````` 
#### Objeto `products`

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| sku | string |  | SKU do produto |
| name | string |  | Nome do produto |
| quantity | number |  | Quantidade do produto |
| unit_value | number |  | Valor unitário do produto (em centavos) |
| type | string |  | Tipo do produto (physical ou digital, padrão: physical) |
```````````````` 
## Exemplo de Resposta (JSON)

#### `201` – Pedido criado com sucesso (retorna os dados do pedido no body)

```
{
  "data": {
    "order": {
      "id": 1,
      "status": "aprovado"
    }
  }
}
```

#### `404` – Cliente ou dados necessários não encontrados

```
{
  "error": {
    "message": "Merchant not found"
  }
}
```

#### `422` – Erro ao validar os dados do pedido

```
{
  "errors": {
    "message": {
      "products_value": [
        "The products value must be an integer."
      ],
      "shipping_value": [
        "The shipping value must be an integer."
      ],
      "discount_value": [
        "The discount value must be an integer."
      ],
      "products": [
        "The products field is required."
      ]
    }
  }
}
```

> 
### 
Dica

Armazene o `order_id` retornado nesta etapa, pois ele será necessário para **efetuar o pagamento** ou **consultar o status do pedido** nas próximas etapas.

---

<!-- url: https://appmax.readme.io/reference/c%C3%A1lculo-do-valor-total-da-order -->
## 4.1.2. Cálculo do valor total do Pedido

### O cálculo do valor total de uma order segue a seguinte regra:

1. **Cálculo Baseado no Unit Value dos Produtos:** 
  - Se o valor unitário (`unit_value`) de cada produto for informado, o valor total será a soma do valor de todos os produtos. Por exemplo, se forem passados três produtos com valores de R$ 10,00, R$ 30,00 e R$ 50,00, o valor total dos produtos será R$ 90,00.
  - Caso a compra seja parcelada, os juros devem ser calculados sobre o valor total dos produtos somado ao frete (`shipping_value`). Por exemplo, se o valor dos produtos for de R$ 90,00 e o frete for R$ 15,00, o valor total a ser considerado para o cálculo dos juros será R$ 105,00. Se for aplicado um juros de 10% em 5x, o valor final será de R$ 115,50 (R$ 105,00 + R$ 10,50 de juros). Esse valor ajustado deve ser distribuído entre os produtos proporcionalmente.
 

### Exemplo de Requisição (cURL)

  -  
```
curl --request POST \
     --url https://api.appmax.com.br/v1/orders \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
    "customer_id": 113543689,
    "discount_value": 0,
    "shipping_value": 1500,
    "products": [
        {
            "sku": "46_0",
            "name": "PRODUTO_TEST_PARA_2",
            "quantity": 1,
            "unit_value": 1350
        },
        {
            "sku": "47_0",
            "name": "PRODUTO_TEST_PARA_2",
            "quantity": 1,
            "unit_value": 3350
        },
        {
            "sku": "48_0",
            "name": "PRODUTO_TEST_PARA_2",
            "quantity": 1,
            "unit_value": 5350
        }
    ]
}
```

2. **Cálculo Baseado no Products Value com Juros:** 
  - Se o valor unitário (`unit_value`) de cada produto não for informado, é obrigatório passar o valor total dos produtos (`products_value`) com juros e frete já calculados. Ou seja, a soma do valor total dos produtos deve incluir tanto os juros quanto o frete. Por exemplo, se o valor total dos produtos e frete for R$ 105,00 e os juros forem de 10%, o valor final será de R$ 115,50. Esse valor deve ser informado diretamente.
 

### Exemplo de Requisição (cURL)

-  
```
curl --request POST \
     --url https://api.appmax.com.br/v1/orders \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
    "customer_id": 113543689,
    "products_value": 10050,
    "discount_value": 0,
    "shipping_value": 1500,
    "products": [
        {
            "sku": "46_0",
            "name": "PRODUTO_TEST_PARA_2",
            "quantity": 1            
        },
        {
            "sku": "47_0",
            "name": "PRODUTO_TEST_PARA_2",
            "quantity": 1            
        },
        {
            "sku": "48_0",
            "name": "PRODUTO_TEST_PARA_2",
            "quantity": 1            
        }
    ]
}
```

### Regras Gerais

- Sempre é necessário enviar o cálculo com os juros incluídos, seja no valor de cada produto individualmente ou no valor total dos produtos e frete.
- O sistema não faz o cálculo dos juros automaticamente; o valor informado já deve estar ajustado de acordo com a forma de pagamento, os juros aplicáveis e o frete. Se os juros forem aplicados, o valor final deverá refletir esse acréscimo, considerando também o frete no cálculo.

---

<!-- url: https://appmax.readme.io/reference/consultar-dados-de-um-pedido -->
## 4.2. Consultar Dados de um Pedido

Este endpoint permite consultar os detalhes de um pedido previamente criado para um merchant. Basta informar o ID do pedido na URL.

## 1. Autenticação

A requisição precisa incluir um **token Bearer**, obtido no fluxo de autenticação do **merchant**.

#### Headers

| Nome | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| Authorization | string |  | Token de autenticação no formato Bearer {TOKEN}. |
`````` 
## 2. Exemplo de Requisição (cURL)

```
curl --location 'https://api.appmax.com.br/v1/orders/3531' \
--header 'Authorization: Bearer SEU_TOKEN'
```

## 3. Exemplo de Resposta (JSON)

```
{
    "data": {
        "order": {
            "id": 3531,
            "status": "estornado",
            "total_paid": 8916,
            "amounts": {
                "sub_total": 4662,
                "shipping_value": 2738,
                "discount": 0,
                "installment_fee": 1516
            },
            "created_at": "2025-02-13 14:09:48",
            "updated_at": "2025-02-13 14:11:55"
        },
        "customer": {
            "id": 2023,
            "name": "Junior Almeida",
            "email": "[email protected]",
            "document_number": "19100000000"
        },
        "payment": {
            "method": "creditcard",
            "installments": 12,
            "installments_amount": 743,
            "card": {
                "brand": "visa",
                "number": "400000****0010"
            },
            "paid_at": "2025-02-13 14:10:20"
        },
        "refund": {
            "refunded_at": "2025-02-13 14:11:55"
        }
    }
}
```

### Observações

- O objeto `payment` pode estar vazio se o pagamento ainda não foi processado.
- O objeto `refund` pode estar vazio caso não haja reembolso no pedido.
### Resumo

- Consulta um pedido existente usando apenas o `order_id`.
- Retorna dados completos do pedido, pagamento, reembolso e cliente.
- Requer autenticação com token `Bearer` gerado com as credenciais do **merchant**.

---

<!-- url: https://appmax.readme.io/reference/como-efetuar-um-pagamento-atraves-da-appmax -->
## 5. Como efetuar um pagamento

A API da Appmax permite a criação de pagamentos para pedidos já existentes utilizando diferentes métodos.

### Métodos disponíveis:

- **Apple Pay**
- **Cartão de Crédito**
- **Pix**
- **Boleto**
> 
### 
Pré-requisito geral:

Antes de criar um pagamento, você precisa ter:

- `order_id` → ID do pedido
- `customer_id` → ID do cliente
 Veja como obtê-los nas seções:

- Criar/Atualizar um Cliente
- Criar um Pedido

---

<!-- url: https://appmax.readme.io/reference/appmax-pagamento-apple-pay -->
## 5.1. Apple Pay

### Sobre o Apple Pay

O Apple Pay é uma solução de pagamento digital da Apple que permite aos usuários realizarem compras com cartão de crédito de forma rápida, segura e conveniente por meio de dispositivos compatíveis, como iPhone, iPad, Apple Watch e Mac.

Utilizando tecnologia de tokenização e autenticação biométrica (Face ID ou Touch ID), o Apple Pay elimina a necessidade de inserir manualmente dados do cartão em cada transação, reduzindo o tempo de checkout e aumentando a conversão de pagamentos. Além disso, transações realizadas via Apple Pay não possuem risco de receberem chargebacks por fraude.

> 
### 
**Importante**: O Apple Pay é compatível apenas com o navegador Safari e dispositivos Apple com suporte a este método de pagamento.

### Glossário

- **Apple Pay**: Solução de pagamento da Apple que permite compras seguras em sites e apps utilizando o Apple Wallet, com autenticação biométrica e criptografia de ponta a ponta
- **Apple Token (Apple Pay Token)**: Objeto JSON criptografado retornado pela PaymentSheet após o usuário confirmar o pagamento, contendo os campos paymentData, paymentMethod e transactionIdentifier. É usado para efetivar o pagamento pela API Appmax.
- **PaymentSheet**: Interface nativa do Apple Pay que exibe ao cliente o valor, os métodos de pagamento disponíveis e solicita autenticação (Face ID, Touch ID ou senha).
- **Merchant Session da Apple**: Sessão JSON assinada pela Apple, válida para um domínio e merchantIdentifier cadastrados, necessária para validar seu site/app e iniciar fluxos de Apple Pay.
### 1. Instalação do Script

Adicione o script do AppmaxScripts ao seu projeto inserindo o seguinte código dentro da seção ou antes do fechamento do :

#### Pré-requisitos

- Navegador compatível
 O botão do Apple Pay só será exibido no Safari (macOS ou iOS). Em outros browsers, não aparecerá. 
- Carteira Apple Pay configurada
 É necessário ter pelo menos um cartão Visa ou Mastercard adicionado no Apple Wallet.
#### Funcionamento básico

O fluxo de pagamento com Apple Pay via Appmax envolve três componentes principais:

1. Script JS minificado: <https://scripts.appmax.com.br/appmax.min.js>
2. Domínio validado junto à Apple: Registrado via merchant session e via arquivo de associação.
3. Requisição de pagamento: Enviada à API Appmax contendo: 
  1. appleToken válido (gerado pelo PaymentSheet)
  2. order_id válido
  3. customer_id válido
Os próximos passos irão mostrar como ter acesso a todos esses componentes e fazê-los trabalhar em conjunto

### 2. Autorizar a instalação do aplicativo

Antes de exibir o botão, registre seu aplicativo na Appmax e informe o domínio:

```
curl --location 'https://api.appmax.com.br/app/authorize'   
  --header 'Content-Type: application/json'   
  --header 'Authorization: Bearer SEU_TOKEN'   
  --data '{  
    "app_id": "APP_ID",  
    "external_key": "EXTERNAL_KEY",  
    "url_callback": "URL_CALLBACK",  
    "domain_name": "subdominio.dominio.com.br"  
}'
```

#### Observação:

O parâmetro `domain_name` deve conter subdomínio + domínio (ex.: `minhaloja.minhaintegracao.com.br`). 

Mais detalhes em Autorização de Aplicativo.

### 3. Disponibilizar o arquivo de validação de domínio

Cada loja precisa servir, na raiz do seu site, o arquivo:

`.well-known/apple-developer-merchantid-domain-association`

**Exemplo de URL completa:**: `https://subdominio.dominio.com.br/.well-known/apple-developer-merchantid-domain-association`

>  
Veja Configuração de Domínios para Apple Pay via Appmax.

### 4. Configuração e adição do script no front-end

> 
### 
Atenção: Esperamos o `external_id` que você utilizou no retorno HTTP Code 200, em Instalação do aplicativo

O script `appmax.min.js` faz três coisas principais:

1. Estiliza o botão com o design oficial da Apple (opcional). 
2. Inicializa o botão para abrir uma PaymentSheet com os dados do carrinho. 
3. Retorna callbacks de sucesso, erro e atualização quando o usuário interage.
#### 4.1 Inicializar o script com onSuccess, onError , externalId, onUpdate e onAuthorize

```
<script>
  // Callback de sucesso: recebe o Apple Token
  const onAutorize = async (appleToken) => {
    await processarPagamento(appleToken);
  };

  // Callback de erro
  const onError = (err) => {
    console.error('Erro no Apple Pay:', err);
  };
  
  const onSuccess = (data) => {
    customer.value.ip = data.ip || 'IP não encontrado.'
  }
  
  // Callback de update: deve retornar os dados atuais do checkout
  const onUpdate = () => getCheckoutData();

  onMounted(() => {
    const script = document.createElement('script');
    script.src = 'https://scripts.sandboxappmax.com.br/appmax.min.js';
    script.onload = () => {
      if (window.AppmaxScripts) {
        // Passe os parâmetros na ordem correta:
        window.AppmaxScripts.init(onSuccess, onError, externalId, onUpdate, onAuthorize);
      } else {
        console.error('AppmaxScripts não carregado.');
      }
    };
    document.head.appendChild(script);
  });
</script>
```

#### 4.2. Mantendo onUpdate sempre atualizado (JS Puro)

```
// 1. Captura referências aos campos:  
const installmentsInput = document.getElementById('installments');  
const freightInput      = document.getElementById('freight');  
const discountInput     = document.getElementById('discount');  
const productItems      = () => document.querySelectorAll('.product-item');

// 2. Monta o objeto de checkout:  
function getCheckoutData() {  
  const products = Array.from(productItems()).map(item => ({  
    name:     item.querySelector('.product-name').textContent,  
    price:    parseFloat(item.querySelector('.product-price').value),  
    quantity: parseInt(item.querySelector('.product-quantity').value, 10)  
  }));  
  const freight      = parseFloat(freightInput.value  || 0);  
  const discount     = parseFloat(discountInput.value || 0);  
  const totalItems   = products.reduce((sum, p) => sum + p.price \* p.quantity, 0);  
  return {  
    orderId:      sessionStorage.getItem('order_id') || '',  
    total:        totalItems + freight - discount,  
    freight:      freight,  
    discount:     discount,  
    installments: parseInt(installmentsInput.value, 10) || 1,  
    products:     products  
  };  
}
```

#### 4.3. Processar o pagamento

```
async function processarPagamento(appleToken) {  
  // 1. Criar customer (POST /v1/customers)  
  const customerPayload = {  
    first_name: "João",  
    last_name: "Silva",  
    email: "[[email protected]](mailto:[email protected])",  
    phone: "51985081457",  
    document_number: "34545284027"  
  };  
  const customerResp = await fetch('https://api.appmax.com.br/v1/customers', {  
    method:  'POST',  
    headers: {  
      'Content-Type': 'application/json',  
      Authorization:  `Bearer SEU_TOKEN`  
    },  
    body: JSON.stringify(customerPayload)  
  });  
  if (!customerResp.ok) throw new Error('Falha ao criar customer');  
  const { id: customer_id } = await customerResp.json();

  // 2. Criar order (POST /v1/orders)  
  const orderPayload = {  
    customer_id,  
    products_value: 6178,  
    discount_value: 0,  
    shipping_value: 2738,  
    products: [  
      {  
        sku:  "758",  
        name: "TÊNIS CASUAL MASCULINO CNS BOLD I PRETO",  
        quantity: 1  
      }  
    ]  
  };  
  const orderResp = await fetch('https://api.appmax.com.br/v1/orders', {  
    method:  'POST',  
    headers: {  
      'Content-Type': 'application/json',  
      Authorization:  `Bearer SEU_TOKEN`  
    },  
    body: JSON.stringify(orderPayload)  
  });  
  if (!orderResp.ok) throw new Error('Falha ao criar order');  
  const { id: order_id } = await orderResp.json();

  // 3. Efetivar pagamento Apple Pay (POST /v1/payments/apple-pay)  
  const paymentPayload = {  
    payment_data: {  
      apple_pay: {  
        installments: "3",  
        holder_document_number: "22233344450",  
        soft_descriptor: "EXEMPLOLOJA",  
        payment_data: {  
          data: "exemplo",  
          signature: "signature",  
          header: {  
            publicKeyHash: "Z...",  
            ephemeralPublicKey: "MFk...==",  
            transactionId: "trx....wvu"  
          },  
          version: "EC_v1"  
        },  
        payment_method: {  
          displayName: "Visa •••• 3714",  
          network:     "Visa",  
          type:        "credit"  
        },  
        transaction_identifier: "trx....wvu"  
      }  
    },  
    order_id,  
    customer_id  
  };  
  const paymentResp = await fetch('https://api.appmax.com.br/v1/payments/apple-pay', {  
    method:  'POST',  
    headers: {  
      'Content-Type': 'application/json',  
      Authorization:  `Bearer SEU_TOKEN`  
    },  
    body: JSON.stringify(paymentPayload)  
  });  
  if (!paymentResp.ok) throw new Error('Falha no pagamento');  
  return paymentResp.json();  
}
```

Para que a chamada à API Appmax de Apple Pay funcione com sucesso, você precisa reunir três informações principais:

- Customer ID válido 
- Identificador do cliente. 
Saiba como obtê-lo em: API de Customers  

- Order ID válido 
- Identificador do pedido. Saiba como obtê-lo em: API de Orders 
- Apple Token: Objeto criptografado retornado pela PaymentSheet no front-end, contendo paymentData, paymentMethod e transactionIdentifier. 
Saiba como obtê-lo em: API de Apple Pay 

Em seguida, basta enviar um `POST` para o endpoint `/v1/payments/apple-pay` incluindo esses três elementos no corpo da requisição seguindo o payload especificado na documentação API de Apple Pay 

### Diagrama de sequência do fluxo

> [Image: # API Payment Flow Diagram

This is a **sequence diagram** illustrating the complete Apple Pay payment flow across four actors: User, Client, Backend, and API Appmax. It shows 10 sequential steps, starting with user authorization and PaymentSheet opening (step 1), progressing through payment function calls with Apple Token (step 2), customer/order creation (steps 3-6), payment execution via the Appmax API (steps 7-8), and concluding with status confirmation returned to the user (steps 9-10).]

### Exemplo de Apple Token

Para fins de verificação, a estrutura de um token Apple deve se assemelhar a um objeto como este:

```
{  
  "paymentData": {  
    "version": "EC_v1",  
    "data": "eyJ2ZXJzaW9uIjoiRUM...",  
    "signature": "MEYCIQDn6C1u...",  
    "header": {  
      "publicKeyHash": "xQiDLOvxMW/xPEIjogYSZBcwkn8uCpgS77ykmI5406g=",  
      "ephemeralPublicKey": "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEDYIDn6Wu/HbEL/cHaU9n41iwJig1GPbowphe5TBTmZwEKv4u0ntW6E7dlc94iG3wmbT3XvpkhdpHkbaPQXk8aA==",  
      "transactionId": "af1252480b5d7f4cfd12e17f45df898bb816dec225b3d55380b9339983044978"  
    }  
  },  
  "paymentMethod": {  
    "displayName": "Visa •••• 4242",  
    "network": "Visa",  
    "type": "debit"  
  },  
  "transactionIdentifier": "000000-123456-7890"  
}
```

### Conclusão

Seguindo este guia, seu checkout exibirá o botão Apple Pay no Safari e permitirá pagamentos com um clique, de forma segura e integrada à API Appmax.

**Lembre-se de:**

- Autorizar a instalação com `domain_name`. 
- Validar o domínio com o arquivo `.well-known/....`
- Incluir e inicializar o script da Appmax (passando onSuccess, onError e onUpdate). 
- Processar o appleToken via API Appmax

---

<!-- url: https://appmax.readme.io/reference/appmax-pagamento-cartao-de-credito -->
## 5.2. Pagamento com Cartão de Crédito

Entenda como o que é necessário e como efetuar um pagamento com cartão de crédito na Appmax

O pagamento por cartão de crédito permite parcelamento, personalização do soft descriptor (texto exibido na fatura) e integração direta ao fluxo de checkout.

Existem duas formas de processar pagamentos com cartão de crédito:

- **Via Appmax JS (CDN):** captura os dados do cartão diretamente no navegador do cliente usando um JavaScript fornecido pela Appmax, garantindo conformidade com o PCI DSS sem que os dados sensíveis passem pelo seu servidor.
- **Via Tokenização (API):** substitui os dados reais do cartão (número, CVV e validade) por um token seguro, permitindo realizar transações sem expor as informações originais.
Esta rota será utilizada para informar os valores com taxas de parcelamento definidas na configuração de pagamento do site cadastradas na Appmax. Ela retorna o valor de cada parcela com base no total do pedido e no tipo de formatação selecionado. Existem duas modalidades disponíveis:

**PP - “Simples por Parcela”:** Nesta modalidade, a taxa de juros é aplicada diretamente sobre o valor de cada parcela. O custo adicional é somado ao valor base da parcela, proporcionando uma visualização clara e direta do valor final a ser pago em cada mês.
 **AM - “Financiamento”:** Nesta modalidade, a taxa de juros é calculada mensalmente sobre o saldo devedor total. O valor das parcelas varia, considerando a aplicação de juros sobre o saldo total a cada mês.
 Essas modalidades atendem a diferentes necessidades de pagamento, sendo o "PP" a mais utilizada, enquanto o "AM" é configurado em situações específicas, mas igualmente relevante. Por isso, é imprescindível consultar a rota para garantir que os valores dos pedidos sejam processados de forma consistente em ambos os sistemas, utilizando a mesma taxa.

Essa integração também facilita a usabilidade, eliminando a necessidade de configurar manualmente as taxas de parcelamento em sistemas externos, já que os valores estarão alinhados com as configurações definidas na Appmax.

#### Observações:

A quantidade e os valores das parcelas podem ser personalizados de 1 a 12 parcelas.
 Os valores são configuráveis individualmente para cada merchant.

### Exemplo de Requisição (cURL)

`POST /v1/payments/credit-card` — Cria um novo pagamento por cartão de crédito vinculado a um pedido existente.

### Exemplo de requisição:

```
curl --request POST \
     --url https://api.appmax.com.br/v1/payments/credit-card \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '{
  "order_id": 12345,
  "customer_id": 407,
  "payment_data": {
    "credit_card": {
      "token": "422146c7523a46119d6073ea56193913",
      "upsell_hash": "12535191020163-1573114925-0097661001572114921",
      "number": "4444222222222222",
      "cvv": "123",
      "expiration_month": "1",
      "expiration_year": "25",
      "holder_document_number": "19100000000",
      "holder_name": "Teste teste",
      "installments": 1,
      "soft_descriptor": "MYSTORE"
    }
  }
}'

```

### Parâmetros do body:

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| order_id | number |  | ID do pedido |
| customer_id | number |  | ID do cliente |
| payment_data | object |  | Detalhes do pagamento |
| credit_card | object |  | Dados do cartão de crédito |
```````` 
### Respostas:

- `201` – Pagamento criado com sucesso
- `400` – Erro ao processar (ex.: pedido já pago)
- `404` – Pedido não encontrado

---

<!-- url: https://appmax.readme.io/reference/421c%C3%A1lculo-de-parcelas -->
## 5.2.1. Cálculo de parcelas

Esta rota será utilizada para informar os valores já com as taxas de parcelamento definidas na configuração de pagamento do site cadastradas na Appmax.

Retorna o valor total com os juros aplicados em cada modalidade de parcelamento. A partir desse valor, a integração deverá realizar a divisão do seu lado para obter o valor exato de cada parcela, de acordo com a quantidade selecionada e o tipo de formatação configurado. Existem duas modalidades disponíveis:

PP - “Simples por Parcela”: Nesta modalidade, a taxa de juros é aplicada diretamente sobre o valor de cada parcela. O custo adicional é somado ao valor base da parcela, proporcionando uma visualização clara e direta do valor final a ser pago em cada mês.
 AM - “Financiamento”: Nesta modalidade, a taxa de juros é calculada mensalmente sobre o saldo devedor total. O valor das parcelas varia, considerando a aplicação de juros sobre o saldo total a cada mês.
 Essas modalidades atendem a diferentes necessidades de pagamento, sendo o "PP" a mais utilizada, enquanto o "AM" é configurado em situações específicas, mas igualmente relevante. Por isso, é imprescindível consultar a rota para garantir que os valores dos pedidos sejam processados de forma consistente em ambos os sistemas, utilizando a mesma taxa.

Essa integração também facilita a usabilidade, eliminando a necessidade de configurar manualmente as taxas de parcelamento em sistemas externos, já que os valores estarão alinhados com as configurações definidas na Appmax.

Observações:

A quantidade e os valores das parcelas podem ser personalizados de 1 a 12 parcelas.
 Os valores são configuráveis individualmente para cada merchant.

Para saber mais sobre como calcular o valor total da Order e garantir o envio correto dos valores na criação de pedidos, clique aqui.

```
curl --request POST \
     --url https://api.appmax.com.br/v1/payments/installments \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
  "installments": 10,
  "total_value": 10000,
  "settings": true
}
```

---

<!-- url: https://appmax.readme.io/reference/tokenizacao-do-cartao-de-credito -->
## 5.2.2. Tokenização do Cartão de Crédito via API

POST: https://api.appmax.com.br/v1/payments/tokenize

A tokenização é realizada após o envio dos dados do cartão (como número do cartão, CVV, etc.) para a API. 

#### Exemplo:

```
curl --request POST \
     --url https://api.appmax.com.br/v1/payments/tokenize \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
  "payment_data": {
    "credit_card": {
      "number": "4444222222222222",
      "cvv": "123",
      "expiration_month": "12",
      "expiration_year": "28",
      "holder_name": "John Doe"
    }
  }
}
'
```

---

<!-- url: https://appmax.readme.io/reference/tokenizacao-com-appmax-js -->
## 5.2.3. Tokenização de Pagamento com Appmax JS

A tokenização é realizada quando um formulário de pagamento é enviado com o atributo data-appmax-checkout. Durante este processo, os dados sensíveis do cartão (como número do cartão, CVV, etc.) são coletados e convertidos em um token seguro. O token é então enviado ao servidor em vez dos dados reais, garantindo que informações sensíveis não sejam expostas.

> 
### 
Esperamos o external_id que você utilizou no retorno HTTP Code 200, em Instalação do aplicativo

#### Exemplo de Formulário de Pagamento:

```
<form id="payment-form" method="POST" data-appmax-checkout>
      <div class="form-group">
        <label for="card-number">Número do Cartão:</label>
        <input type="text" id="card-number" name="card-number" appmax-form-element="number" required>
      </div>

      <div class="form-group">
        <label for="card-holder-name">Nome do Titular:</label>
        <input type="text" id="card-holder-name" name="card-holder-name" appmax-form-element="holder_name" required>
      </div>

      <div class="form-group">
        <label for="exp-month">Mês de Expiração:</label>
        <input type="text" id="exp-month" name="exp-month" appmax-form-element="expiration_month" required>
      </div>

      <div class="form-group">
        <label for="exp-year">Ano de Expiração:</label>
        <input type="text" id="exp-year" name="exp-year" appmax-form-element="expiration_year" required>
      </div>

      <div class="form-group">
        <label for="cvv">CVV:</label>
        <input type="text" id="cvv" name="cvv" appmax-form-element="cvv" required>
      </div>
  
      <!-- O token será gerado e enviado ao servidor no lugar dos dados reais -->

      <button type="submit">Pagar</button>
</form>
```

#### Parâmetros do Formulário:

Os dados dos cartões de crédito são coletados a partir dos seguintes campos utilizando nos inputs o `appmax-form-element`:

- `number:` Número do cartão de crédito.
- `holder_name:` Nome do titular do cartão.
- `expiration_month: `Mês de expiração do cartão.
- `expiration_year:` Ano de expiração do cartão.
- `cvv:` Código de segurança do cartão (CVV).

---

<!-- url: https://appmax.readme.io/reference/appmax-pagamento-pix -->
## 5.3. Pagamento via Pix

Entenda como o que é necessário e como efetuar um pagamento com Pix na Appmax

O Pix oferece **pagamento instantâneo**, disponível 24/7, com confirmação muito rápida.

Ao integrar este endpoint, você pode gerar um **QR Code** e/ou o **código EMV** para exibir ao cliente no checkout, possibilitando a conclusão do pagamento em poucos segundos.

Importante: exiba um cronômetro na página de sucesso junto com o QR Code e a opção de copiar a chave, para que o cliente saiba o tempo restante de expiração da chave.

> 
### 
Importante: exiba um cronômetro na página de sucesso junto com o QR Code e a opção de copiar a chave, para que o cliente saiba o tempo restante de expiração da chave.

### Endpoint:

`POST /v1/payments/pix` — Gera as instruções de pagamento via Pix para um pedido existente.

### Exemplo de Requisição (cURL)

```
curl --request POST \
     --url https://api.appmax.com.br/v1/payments/pix \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
  "order_id": 113,
  "payment_data": {
    "pix": {
      "document_number": "19100000000"
    }
  }
}
'
```

### Parâmetros do body:

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| order_id | number |  | ID do pedido |
| payment_data | object |  | Objeto com dados do pagamento |
| pix | object |  | Objeto vazio para indicar método Pix |
`````` 
### Respostas:

- `200` – Pagamento criado com sucesso (retorna QR Code e EMV)
- `404` – Pedido já pago ou não encontrado

---

<!-- url: https://appmax.readme.io/reference/appmax-pagamento-boleto -->
## 5.4. Pagamento via Boleto

Entenda como o que é necessário e como gerar um pagamento com Boleto na Appmax

O boleto bancário é uma opção prática para clientes que preferem ou precisam pagar de forma offline, em bancos, lotéricas ou aplicativos de internet banking.

> 
### 
Importante: exiba na página de sucesso a opção de baixar o boleto e copiar a linha digitável, garantindo que o cliente tenha acesso fácil ao pagamento.

Obs.: Não pode ser exibido via Iframe, precisa redirecionar para outra pagina ao clicar em baixar boleto.

### Endpoint:

`POST /v1/payments/boleto` — Cria um boleto vinculado a um pedido existente.

### Exemplo de Requisição (cURL)

```
curl --request POST \
     --url https://api.appmax.com.br/v1/payments/boleto \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '{
  "order_id": 113,
  "payment_data": {
    "boleto": {
      "document_number": "19100000000"
    }
  }
}'
```

### Parâmetros do body:

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| order_id | number |  | ID do pedido |
| payment_data | object |  | Detalhes do pagamento |
| boleto | object |  | Dados do boleto |
| document_number | string |  | CPF ou CNPJ do pagador |
```````` 
### Respostas:

- `201` – Pagamento criado com sucesso (retorna link PDF e linha digitável)
- `404` – Pedido já pago ou não encontrado

---

<!-- url: https://appmax.readme.io/reference/estorno-na-appmax -->
## 6. Criar um Estorno

Entenda como fazer o estorno de um pedido na Appmax

O endpoint `POST /v1/orders/refund-request` é utilizado para criar um **estorno** na Appmax.

 Um estorno deve estar sempre vinculado a um **pedido previamente criado** no sistema.

> 
### 
Pré-requisito:

Para criar um estorno, você** precisa ter** o `order_id` do pedido.

 Caso ainda não tenha criado o cliente, siga antes a documentação: 

Como criar um Pedido na Appmax

### Exemplo de Requisição (cURL)

```
curl --request POST \
     --url https://api.appmax.com.br/v1/orders/refund-request \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
  "order_id": 1,
  "type": "total",
  "value": 10000
}
'
```

### Parâmetros do Body:

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| order_id | number |  | ID do Pedido (obtido na etapa de criação do pedido) |
| type | string |  | Tipo de reembolso (total ou partial, padrão: total) |
| value | string | ‼ | Valor do reembolso |
```````````` 
> 
### ℹ
O `value` é obrigatório caso o `type` seja `partial`

## Exemplo de Resposta (JSON)

#### `201` – Estorno criado com sucesso

```
{
  "data": {
    "message": "Refund request accepted"
  }
}
```

#### `400` – Erro de validação

```
{
  "error": {
    "message": "Error to refund"
  }
}
```

#### `404` – Pedido não encontrado

```
{
  "error": {
    "message": [
      "Order not found",
      "Site not found"
    ]
  }
}
```

#### `500` – Erro ao genérico

```
{
  "error": {
    "message": "Error to refund"
  }
}
```

---

<!-- url: https://appmax.readme.io/reference/criar-um-upsell-na-appmax -->
## 7. Criar um Upsell

Como criar um Upsell de um Pedido na Appmax.

O endpoint `POST /v1/orders/upsell` é utilizado para criar um **upsell** na Appmax.

 Um upsell deve estar sempre vinculado a um **pedido** e um **pagamento** previamente criado no sistema.

> 
### 
Pré-requisito:

Para criar um estorno, você** precisa ter** o `upsell_hash` gerado no pagamento do pedido.

 Caso ainda não tenha criado o cliente, siga antes a documentação: 

Como criar um Pedido na Appmax

### Exemplo de Requisição (cURL)

```
curl --request POST \
     --url https://api.appmax.com.br/v1/orders/upsell \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
  "upsell_hash": "4000114202503117156088040208561001715608804",
  "products_value": 10000,
  "products": [
    {
      "sku": "9000010",
      "name": "Livro de receitas",
      "quantity": 1,
      "unit_value": 12300,
      "type": "digital"
    }
  ]
}
'
```

### Parâmetros do Body:

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| upsell_hash | string |  | Hash do pedido para Upsell |
| products_value | number |  | Valor total dos produtos do Upsell |
| products | array |  | Produtos do Upsell |
`````` 
#### Objeto `products`

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| sku | string |  | SKU do produto |
| name | string |  | Nome do produto |
| quantity | number |  | Quantidade do produto |
| unit_value | number |  | Valor unitário do produto (em centavos) |
| type | string |  | Tipo do produto (physical ou digital, padrão: physical) |
```````````````` 
## Exemplo de Resposta (JSON)

#### `201` – Upsell criado com sucesso

```
{
  "data": {
    "message": "Transação efetuada com sucesso",
    "redirect_url": "example.com/order/success-by-order?hash=40001142025031-1715608804-0190017001715608804"
  }
}
```

#### `404` – Pedido não encontrado

```
{
  "data": {
    "errors": {
      "message": "Order not found"
    }
  }
}
```

#### `422` – Erro ao validar os dados

```
{
  "errors": {
    "message": {
      "upsell_hash": [
        "The upsell hash field is required.",
        "The Upsell Hash must be a string."
      ]
    }
  }
}
```

---

<!-- url: https://appmax.readme.io/reference/como-criar-recorrencia-na-appmax -->
## 8. Recorrência na Appmax

Como criar um pedido recorrente na Appmax

A recorrência é uma funcionalidade oferecida pela Appmax que permite ao merchant configurar cobranças periódicas.
 Para utilizá-la, o merchant precisa:

- Instalar um aplicativo de recorrência internamente em seu acesso na Appmax.
- Configurar a recorrência manualmente em sua plataforma.
Para que a recorrência funcione no seu aplicativo, é necessário enviar, nos endpoints de **cartão de crédito e Pix**, as informações do objeto `subscription`, contendo os seguintes campos:

- **interval:** define o período da recorrência, podendo ser mensal, anual ou semanal.
- **interval_count:** representa a quantidade de ciclos do intervalo definido, aceitando valores de 1 a 12.
O objeto `payment_data` do cartão e do Pix permanece o mesmo, basta adicionar o objeto subscription. Essa configuração funciona tanto para cartões tokenizados quanto para cartões não tokenizados.

```
{
    "order_id": 99,
    "payment_data": {
        //credit_card, pix..
        "subscription": {
            "interval": "month",
            "interval_count": 1
        }
    }
}
```

#### Parâmetros do Body:

| Nome | Tipo | Obrigatório | Valores Validos |
| --- | --- | --- | --- |
| interval | string |  | month ou year |
| interval_count | string |  | valores de 1 a 12. |
````````````

---

<!-- url: https://appmax.readme.io/reference/cadastrar-codigo-de-rastreio-na-appmax -->
## 9. Como cadastrar o código de rastreio no Pedido

Entenda como cadastrar o código de rastreio dos pedidos na Appmax

O endpoint `POST /v1/orders/shipping-tracking-code` é utilizado para cadastrar um **código de rastreio** no pedido criado na Appmax.

 Um **código de rastreio** deve estar sempre vinculado a um **pedido previamente criado** no sistema.

> 
### 
Pré-requisito:

Para registrar um **código de rastreio**, você** precisa ter** o `order_id` do pedido.

 Caso ainda não tenha criado o cliente, siga antes a documentação: 

Como criar um Pedido na Appmax

### Exemplo de Requisição (cURL)

```
curl --request POST \
     --url https://api.appmax.com.br/v1/orders/shipping-tracking-code \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
  "order_id": 2,
  "shipping_tracking_code": "EEEASDASDAS1239A"
}
'
```

### Parâmetros do Body:

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| order_id | number |  | ID do Pedido (obtido na etapa de criação do pedido) |
| shipping_tracking_code | string |  | Código de rastreio do pedido |
```` 
## Exemplo de Resposta (JSON)

#### `201` – Código de rastreio incluído com sucesso

```
{
  "data": {
    "message": "tracking accepted"
  }
}
```

#### `400` – Pedido não encontrado

```
{
  "errors": {
    "message": "Failed to store delivery tracking code."
  }
}
```

#### `422` – Erro ao validar os dados

```
{
  "message": "The given data failed to pass validation.",
  "errors": {
    "message": "The delivery_tracking_code field is required."
  }
}
```

---

<!-- url: https://appmax.readme.io/reference/como-utilizar-a-recuperacao-de-vendas-com-ia-da-appmax -->
## 10. Recuperação de Vendas com IA

Entenda como utilizar a Recuperação de Vendas com IA da Appmax

Esta rota permite criar um carrinho abandonado utilizando a **recuperação de vendas com IA** da Appmax.

> 
### 
Pré-requisito:

Para criar um cliente, você** precisa ter** feito a coleta de IP utilizando o script Appmax JS.

 Caso ainda não tenha feito a coleta de IP, siga antes a documentação: 

Coletar IP do Cliente

### Exemplo de Requisição (cURL)

```
curl --request POST \
     --url https://api.appmax.com.br/v1/customers \
     --header 'accept: application/json' \
     --header 'content-type: application/json' \
     --data '
{
  "first_name": "Junior",
  "last_name": "Almeida",
  "email": "[email protected]",
  "phone": "51983655100",
  "document_number": "25226493029",
  "cart_link": "https://subdomain.domain.com/cart/123-345-678",
  "address": {
    "postcode": "91520270",
    "street": "Rua Francisco Carneiro da Rocha",
    "number": "582",
    "complement": "Casa",
    "district": "Moinhos de Ventos",
    "city": "Porto Alegre",
    "state": "RS"
  },
  "ip": "127.0.0.1",
  "products": [
    {
      "sku": "9000010",
      "name": "Livro de receitas",
      "quantity": 1,
      "unit_value": 12300,
      "type": "digital"
    }
  ],
  "tracking": {
    "utm_source": "google",
    "utm_campaign": "teste"
  }
}
'
```

### Parâmetros do Body:

| Campo | Tipo | Obrigatório | Descrição |
| --- | --- | --- | --- |
| first_name | string |  | Nome do cliente |
| last_name | string |  | Sobrenome do cliente |
| email | string |  | E-mail válido |
| phone | string |  | Telefone com DDD |
| document_number | string |  | CPF ou CNPJ |
| address | object |  | Endereço do cliente |
| ip | string |  | IP de origem |
| products | array |  | Lista de produtos vinculados |
| tracking | object |  | Dados de origem da visita |
| cart_link | string |  | URL do carrinho |
```````````````````` 
## Exemplo de Resposta (JSON)

#### `201` – Cliente criado com sucesso

```
{
  "data": {
    "customer": {
      "id": 1
    }
  }
}
```

#### `422` – Erro de validação dos dados

```
{
  "message": "The given data failed to pass validation.",
  "errors": {
    "message": {
      "first_name": [
        "The first_name field is required."
      ],
      "last_name": [
        "The last_name field is required."
      ],
      "phone": [
        "The phone field must be a string.",
        "The phone field must have a maximum of 11 characters."
      ],
      "email": [
        "The email field is required.",
        "The email field must be a string."
      ],
      "ip": [
        "The ip field is required."
      ]
    }
  }
}
```

> 
### 
Importante:

Guarde o valor de `customer_id` retornado nesta etapa, mesmo que temporariamente.

Ele será necessário para criar o pedido na próxima etapa do fluxo.

** Sem este ID, não será possível criar um pedido pois é necessário vincular o cliente ao pedido.**

---

<!-- url: https://appmax.readme.io/reference/faq -->
## FAQ

1.  
**O que fazer quando ocorrer o erro “500 Internal Server Error - an error occurred while integrating the app. Please try again” na autenticação do app? **
 Esse erro geralmente acontece quando o fluxo de instalação do aplicativo não é seguido corretamente, especificamente quando falta o passo de redirecionamento e autorização da instalação.
 O fluxo correto de instalação do app é o seguinte:

  1. **Obter o token do aplicativo**
 Faça um **POST **para:
 https://auth.sandboxappmax.com.br/oauth2/token
 Utilize as credenciais do aplicativo.
 Esse token será usado nas próximas chamadas.
  2. **Gerar o hash de autorização**
 Com o token obtido, faça um **POST **para:
 https://api.sandboxappmax.com.br/app/authorize
 Envie os seguintes parâmetros: 
    1.  
      - app_id
      - client_key
      - url_callback
    2. Esse passo gera um hash de autorização.
  3. **Redirecionar o usuário para autorizar a instalação (passo que geralmente está sendo esquecido)**
 Com o **hash **gerado no passo anterior, redirecione o usuário para:
 https://breakingcode.sandboxappmax.com.br/appstore/integration/HASH
 Esse redirecionamento é essencial para que o **merchant **autorize a instalação do app.
  4. **Gerar as credenciais da API do merchant**
 Após a autorização, use o mesmo **hash **em um **POST **para:
 https://api.sandboxappmax.com.br/app/client/generate
 Esse passo retorna as credenciais **client_id** e **client_secret**, que permitem a utilização da API (para criar pedidos, clientes, pagamentos etc).
  5. **Resumindo:**
 Esse erro ocorre porque está faltando o redirecionamento para o passo de autorização do aplicativo. Sem essa etapa, a instalação não é concluída e as credenciais do merchant não são geradas.
  Documentação oficial do fluxo de autenticação e autorização da Appmax:
 Fluxo de Autenticação e Autorização na API da Appmax
2.  
**Qual o tempo para gerar o token com o ID e secret?** Se essa dúvida for referente ao JWT para acessar a API, ele é um token de curta duração, tem validade de 1h.
 O client_id e client_secret nunca são alterados, só é possível gerar novos realizando novas instalação e desativar os atuais realizando a desinstalação.

3.  
**O que fazer quando ocorrer erro 401 para criar o customer?** 

  1. Verificar se o token está válido, não está expirado.
  2. Verificar se o token utilizado para a requisição foi gerado com as credenciais do merchant após o processo de instalação.
4.  
**Como a integração identifica a loja que está fazendo a instalação do aplicativo?**
 Quando a instalação do aplicativo é iniciada diretamente pela plataforma da integração, a identificação da loja ocorre por meio do login do usuário na plataforma externa.
 Ao clicar em "Instalar" na plataforma da integração, um token deverá ser gerado através da seguinte rota: https://api.appmax.com.br/app/authorize
 Após a geração do token, o merchant será redirecionado para o fluxo de instalação na rota: https://admin.appmax.com.br/appstore/integration/TOKEN_GERADO
 Nesse processo, o merchant deverá informar o nome da loja e selecionar a empresa cadastrada na Appmax (é preciso ter uma conta na Appmax e estar logado nela).
 Quando o merchant é redirecionado para a rota de integração, ele autoriza a instalação do aplicativo. Após essa autorização, o mesmo hash de curta duração deve ser utilizado para gerar as credenciais na seguinte rota: https://api.appmax.com.br/app/client/generate
 **Atenção: **O hash pode ser utilizado apenas uma vez, porém as credenciais (client_id e client_secret) geradas são válidas indefinidamente, até que o aplicativo seja desinstalado.

5.  
**O que fazer quando ocorrer erro 403 na rota /oauth2/token?**
 Nesse caso aqui o erro está ocorrendo porque essa chamada é na rota de autenticação.
 Endpoint correto: https://auth.sandboxappmax.com.br/oauth2/token
 Ou seja, temos o endpoint com a autenticação que é em https://auth.sandboxappmax.com.br e as chamadas da API utilizando o token gerado que aí e em https://api.sandboxappmax.com.br
 **Exemplo do CURL correto:**
 --header 'Content-Type: application/x-www-form-urlencoded' \
 --header 'Cookie: XSRF-TOKEN=26949a8d-9d53-44ff-845c-bd5a09cd1f34' \
 --data-urlencode 'grant_type=client_credentials' \
 --data-urlencode 'client_id=CLIENT_ID' \
 --data-urlencode 'client_secret=CLIENT_SECRET'

```
curl --location 'https://auth.sandboxappmax.com.br/oauth2/token' \

--header 'Content-Type: application/x-www-form-urlencoded' \

--header 'Cookie: XSRF-TOKEN=26949a8d-9d53-44ff-845c-bd5a09cd1f34' \

--data-urlencode 'grant_type=client_credentials' \

--data-urlencode 'client_id=CLIENT_ID' \

--data-urlencode 'client_secret=CLIENT_SECRET'
```

6.  
**Quais os erros de webhook e o significado?** 502 - URL de webhook cadastrado errado

7.  
**Qual a diferença entre external_key e external_id? Pode ser enviado o mesmo?**
 **external_key:**
 É uma chave fornecida pelo sistema da plataforma e serve para identificar a origem da instalação de forma mais precisa. Pode ser algo já existente no ambiente do cliente/lojista (como store_id, merchant_id, company_key etc.). Essa informação não é controlada pela Appmax, apenas recebemos e repassamos para que a plataforma utilize e correlacione com o sistema dela.
 **external_id:**
 É um ID único gerado pelo próprio sistema da plataforma para cada instalação do aplicativo. Ele serve para manter o vínculo interno entre a instalação na Appmax e o seu banco de dados. Essa informação é totalmente controlada pelo backend da própria plataforma.
 **Pode ser o mesmo?**
 Não pode porque cada campo tem uma finalidade distinta:
 **external_key:** identifica a origem no contexto da plataforma/cliente.
 **external_id:** confirmação da instalação do aplicativo.

 **Em resumo:**

| Campo | Quem define | Uso | Pode repetir? | Recomendação |
| --- | --- | --- | --- | --- |
| external_key | Plataforma externa / sistema do lojista | Para a platforma identificar a origem da instalação | Sim | Usar apenas como referência externa |
| external_id | Backend da própria plataforma | Identificar de forma única cada instalação | Não | Sempre gerar internamente e manter exclusivo para cada instalação |
****************************

---
