<script lang="ts">
	import { createDefaultLearnerInput, type NewLearnerInput } from '$lib/modules/learners';
	import LearnerAvatar from '../avatar/LearnerAvatar.svelte';

	let {
		onCreate,
		onCancel
	} = $props<{
		onCreate: (input: NewLearnerInput) => boolean;
		onCancel: () => void;
	}>();

	let form = $state<NewLearnerInput>(createDefaultLearnerInput());
	let isProcessingPhoto = $state(false);
	let photoError = $state('');

	// Reduz a foto antes de salvar no localStorage, evitando imagens pesadas no cadastro.
	function resizeLearnerPhoto(file: File) {
		return new Promise<string>((resolve, reject) => {
			const reader = new FileReader();
			reader.onerror = () => reject(new Error('Nao foi possivel ler a foto.'));
			reader.onload = () => {
				const image = new Image();
				image.onerror = () => reject(new Error('Nao foi possivel carregar a foto.'));
				image.onload = () => {
					const maxSize = 512;
					const scale = Math.min(1, maxSize / Math.max(image.width, image.height));
					const canvas = document.createElement('canvas');
					canvas.width = Math.round(image.width * scale);
					canvas.height = Math.round(image.height * scale);

					const context = canvas.getContext('2d');
					if (!context) {
						reject(new Error('Nao foi possivel processar a foto.'));
						return;
					}

					context.drawImage(image, 0, 0, canvas.width, canvas.height);
					resolve(canvas.toDataURL('image/jpeg', 0.82));
				};
				image.src = String(reader.result ?? '');
			};
			reader.readAsDataURL(file);
		});
	}

	// Atualiza a foto de perfil do aprendente no proprio formulario de criacao.
	async function handlePhotoChange(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;

		isProcessingPhoto = true;
		photoError = '';
		try {
			form.photoUrl = await resizeLearnerPhoto(file);
		} catch (error) {
			photoError = error instanceof Error ? error.message : 'Nao foi possivel processar a foto.';
		} finally {
			isProcessingPhoto = false;
			input.value = '';
		}
	}

	// Envia uma copia do formulario para a pagina orquestradora validar e persistir.
	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		if (onCreate({ ...form })) {
			form = createDefaultLearnerInput();
			isProcessingPhoto = false;
		}
	}
</script>

<form class="add-form" onsubmit={handleSubmit}>
	<!-- Foto do aprendente: aparece em listas, agenda e prontuario depois do cadastro. -->
	<div class="photo-field">
		<LearnerAvatar name={form.name || 'Aprendente'} photoUrl={form.photoUrl} size="large" />
		<div>
			<strong>Foto do aprendente</strong>
			<p>Opcional, mas ajuda a identificar o paciente na agenda e no prontuario.</p>
			<div class="photo-actions">
				<label class="upload-button">
					<input type="file" accept="image/*" onchange={handlePhotoChange} disabled={isProcessingPhoto} />
					<span>{isProcessingPhoto ? 'Processando...' : 'Adicionar foto'}</span>
				</label>
				{#if form.photoUrl}
					<button type="button" class="secondary-button" onclick={() => (form.photoUrl = '')}>
						Remover foto
					</button>
				{/if}
			</div>
			{#if photoError}
				<small>{photoError}</small>
			{/if}
		</div>
	</div>

	<!-- Campos basicos do cadastro: alimentam a criacao do aprendente e sua agenda inicial. -->
	<div class="form-grid">
		<label>
			<span>Nome do aprendente</span>
			<input type="text" bind:value={form.name} required />
		</label>

		<label>
			<span>Genero</span>
			<select bind:value={form.gender}>
				<option value="">Selecionar</option>
				<option value="Feminino">Feminino</option>
				<option value="Masculino">Masculino</option>
				<option value="Outro">Outro</option>
				<option value="Nao informado">Nao informado</option>
			</select>
		</label>

		<label>
			<span>Idade</span>
			<input type="text" bind:value={form.age} placeholder="8 anos" />
		</label>

		<label>
			<span>Responsavel</span>
			<input type="text" bind:value={form.guardian} />
		</label>

		<label>
			<span>Status</span>
			<select bind:value={form.status}>
				<option value="active">Ativo</option>
				<option value="inactive">Inativo</option>
			</select>
		</label>

		<label>
			<span>Inicio</span>
			<input type="date" bind:value={form.startDate} required />
		</label>

		<label>
			<span>Final</span>
			<input type="date" bind:value={form.endDate} required />
		</label>

		<label>
			<span>Numero de visitas</span>
			<input type="number" min="1" max="120" bind:value={form.visitCount} />
		</label>
	</div>

	<!-- Acoes do formulario: cancelar fecha o card, salvar envia para a rota orquestradora. -->
	<div class="form-actions">
		<button type="button" class="secondary-button" onclick={onCancel}>Cancelar</button>
		<button type="submit" class="primary-button">Salvar aprendente</button>
	</div>
</form>
