---
name: commit-push-pr
description: Workflow git do payments-gateway — feature branch → commit → push → PR para main. Invocar com /commit-push-pr quando o trabalho estiver pronto para review.
disable-model-invocation: true
---

Executar o workflow git completo para o payments-gateway:

1. Verificar estado actual:
```bash
git status && git diff
```

2. Se não estiver numa feature branch, criar uma com o formato:
   - `feat/<descrição>` — nova funcionalidade
   - `fix/<descrição>` — correcção de bug
   - `chore/<descrição>` — manutenção, dependências
   - `refactor/<descrição>` — refactorização
   - `test/<descrição>` — testes

3. Stage dos ficheiros relevantes e commit:
   - Mensagem no formato `<tipo>: <descrição concisa em inglês>`
   - **Nunca** incluir `Co-Authored-By`, `Generated with`, ou qualquer referência a AI
   - Usar HEREDOC para garantir formatação correcta

4. Push e criar PR:
```bash
git push -u origin <branch>
gh pr create --title "<tipo>: <título>" --body "$(cat <<'EOF'
## O que muda
- <bullet 1>

## Como testar
- [ ] docker-compose up
- [ ] <passo específico>
EOF
)"
```

5. Retornar o URL do PR.

$ARGUMENTS
