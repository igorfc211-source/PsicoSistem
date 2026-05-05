<script lang="ts">
	let {
		value,
		placeholder = 'Escreva aqui...',
		onChange
	} = $props<{
		value: string;
		placeholder?: string;
		onChange: (value: string) => void;
	}>();

	let editor = $state<HTMLDivElement | null>(null);
	let selectedColor = $state('#1e2430');
	let selectedSize = $state('3');

	$effect(() => {
		if (editor && editor.innerHTML !== value) {
			editor.innerHTML = value;
		}
	});

	// Aplica comandos simples de edicao, criando uma experiencia parecida com um docs leve.
	function applyCommand(command: string, commandValue?: string) {
		editor?.focus();
		document.execCommand(command, false, commandValue);
		emitValue();
	}

	// Mantem o HTML formatado sincronizado com o formulario pai.
	function emitValue() {
		onChange(editor?.innerHTML ?? '');
	}
</script>

<div class="rich-editor">
	<!-- Barra de ferramentas: comandos simples de formatacao para relatorios. -->
	<div class="rich-toolbar" aria-label="Ferramentas de formatacao">
		<button type="button" onclick={() => applyCommand('bold')}>Negrito</button>

		<label>
			<span>Tamanho</span>
			<select
				bind:value={selectedSize}
				onchange={() => applyCommand('fontSize', selectedSize)}
			>
				<option value="2">Pequeno</option>
				<option value="3">Normal</option>
				<option value="5">Grande</option>
				<option value="7">Titulo</option>
			</select>
		</label>

		<label>
			<span>Cor</span>
			<input
				type="color"
				bind:value={selectedColor}
				oninput={() => applyCommand('foreColor', selectedColor)}
			/>
		</label>
	</div>

	<!-- Area editavel: mantem HTML formatado sincronizado com o componente pai. -->
	<div
		bind:this={editor}
		class="rich-editor-area"
		contenteditable="true"
		role="textbox"
		aria-multiline="true"
		aria-label={placeholder}
		data-placeholder={placeholder}
		oninput={emitValue}
	></div>
</div>
