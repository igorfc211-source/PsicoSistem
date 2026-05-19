<script lang="ts">
	import {
		RELATIONSHIP_OPTIONS,
		createEmptyLearnerGuardianInput,
		createDefaultLearnerInput,
		normalizePhone,
		type GuardianOption,
		type NewLearnerInput
	} from '$lib/modules/learners';
	import LearnerAvatar from '../avatar/LearnerAvatar.svelte';

	let {
		onCreate,
		onCancel,
		guardianOptions
	} = $props<{
		onCreate: (input: NewLearnerInput) => boolean;
		onCancel: () => void;
		guardianOptions: GuardianOption[];
	}>();

	let form = $state<NewLearnerInput>(createDefaultLearnerInput());
	let sessionPriceValue = $state('');
	let generalValue = $state('');
	let isProcessingPhoto = $state(false);
	let photoError = $state('');

	const guardianIndexes = [0, 1];

	function parseCurrencyToCents(value: string) {
		const normalizedValue = value.replace(/[^\d,.-]/g, '').replace(',', '.').trim();
		const amount = Number(normalizedValue);
		if (!Number.isFinite(amount)) return 0;

		return Math.max(0, Math.round(amount * 100));
	}

	// Reduz a foto antes de salvar no localStorage, evitando imagens pesadas no cadastro
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

	function handleSessionPriceInput(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		sessionPriceValue = input.value;
		form.sessionPriceCents = parseCurrencyToCents(input.value);
	}

	function handleGeneralValueInput(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		generalValue = input.value;
		form.generalValueCents = parseCurrencyToCents(input.value);
	}

	function getGuardian(index: number) {
		return form.guardians[index] ?? createEmptyLearnerGuardianInput();
	}

	function updateGuardian(index: number, patch: Partial<NewLearnerInput['guardians'][number]>) {
		const guardians = [...form.guardians];
		guardians[index] = {
			...createEmptyLearnerGuardianInput(),
			...getGuardian(index),
			...patch
		};
		form = {
			...form,
			guardians
		};
		syncPrimaryGuardianFields();
	}

	function syncPrimaryGuardianFields() {
		const primaryGuardian = form.guardians.find((guardian) => guardian.name.trim());
		form = {
			...form,
			guardian: primaryGuardian?.name ?? '',
			guardianRelationship: primaryGuardian?.relationship ?? ''
		};
	}

	function handleExistingGuardianSelection(event: Event, index: number) {
		const key = (event.currentTarget as HTMLSelectElement).value;
		const option = guardianOptions.find((item: GuardianOption) => item.key === key);

		if (!option) {
			updateGuardian(index, createEmptyLearnerGuardianInput());
			return;
		}

		updateGuardian(index, {
			sourceKey: option.key,
			name: option.name,
			relationship: option.relationship,
			phone: option.phone
		});
	}

	function handleGuardianNameInput(event: Event, index: number) {
		updateGuardian(index, {
			sourceKey: '',
			name: (event.currentTarget as HTMLInputElement).value
		});
	}

	function handleGuardianRelationshipChange(event: Event, index: number) {
		updateGuardian(index, {
			relationship: (event.currentTarget as HTMLSelectElement).value
		});
	}

	function handleGuardianPhoneInput(event: Event, index: number) {
		updateGuardian(index, {
			phone: normalizePhone((event.currentTarget as HTMLInputElement).value).slice(0, 11)
		});
	}

	function formatPhoneInput(phone: string) {
		const digits = normalizePhone(phone).slice(0, 11);
		if (digits.length <= 2) return digits;
		if (digits.length <= 6) return `(${digits.slice(0, 2)}) ${digits.slice(2)}`;
		if (digits.length <= 10) {
			return `(${digits.slice(0, 2)}) ${digits.slice(2, 6)}-${digits.slice(6)}`;
		}

		return `(${digits.slice(0, 2)}) ${digits.slice(2, 7)}-${digits.slice(7)}`;
	}

	function isGuardianOptionUsed(optionKey: string, currentIndex: number) {
		return form.guardians.some(
			(guardian, index) => index !== currentIndex && guardian.sourceKey === optionKey
		);
	}

	// Envia uma copia do formulario para a pagina orquestradora validar e persistir.
	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		if (onCreate({ ...form })) {
			form = createDefaultLearnerInput();
			sessionPriceValue = '';
			generalValue = '';
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

		{#each guardianIndexes as guardianIndex}
			{@const guardian = getGuardian(guardianIndex)}
			<section class="guardian-form-section">
				<div class="guardian-form-head">
					<strong>Responsavel {guardianIndex + 1}</strong>
					<span>{guardianIndex === 0 ? 'Obrigatorio' : 'Opcional'}</span>
				</div>

				<div class="guardian-form-grid">
					<label>
						<span>Interligar responsavel existente</span>
						<select
							value={guardian.sourceKey}
							onchange={(event) => handleExistingGuardianSelection(event, guardianIndex)}
						>
							<option value="">Cadastrar novo</option>
							{#each guardianOptions as option}
								<option
									value={option.key}
									disabled={isGuardianOptionUsed(option.key, guardianIndex)}
								>
									{option.name}
								</option>
							{/each}
						</select>
					</label>

					<label>
						<span>Nome do responsavel</span>
						<input
							type="text"
							value={guardian.name}
							required={guardianIndex === 0}
							placeholder="Marina Silva"
							oninput={(event) => handleGuardianNameInput(event, guardianIndex)}
						/>
					</label>

					<label>
						<span>Parentesco</span>
						<select
							value={guardian.relationship}
							required={guardianIndex === 0 || Boolean(guardian.name.trim())}
							onchange={(event) => handleGuardianRelationshipChange(event, guardianIndex)}
						>
							<option value="">Selecionar</option>
							{#each RELATIONSHIP_OPTIONS as option}
								<option value={option}>{option}</option>
							{/each}
						</select>
					</label>

					<label>
						<span>Numero</span>
						<input
							value={formatPhoneInput(guardian.phone)}
							inputmode="numeric"
							pattern={'\\([0-9]{2}\\) [0-9]{4,5}-[0-9]{4}'}
							maxlength="15"
							placeholder="(11) 99999-9999"
							title="Informe um numero no formato (xx) xxxxx-xxxx."
							oninput={(event) => handleGuardianPhoneInput(event, guardianIndex)}
						/>
					</label>
				</div>
			</section>
		{/each}

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

		<label>
			<span>Valor por sessao</span>
			<input
				type="text"
				inputmode="decimal"
				value={sessionPriceValue}
				placeholder="150,00"
				oninput={handleSessionPriceInput}
			/>
		</label>

		<label>
			<span>Valor geral</span>
			<input
				type="text"
				inputmode="decimal"
				value={generalValue}
				placeholder="1200,00"
				oninput={handleGeneralValueInput}
			/>
		</label>
	</div>

	<!-- Acoes do formulario: cancelar fecha o card, salvar envia para a rota orquestradora. -->
	<div class="form-actions">
		<button type="button" class="secondary-button" onclick={onCancel}>Cancelar</button>
		<button type="submit" class="primary-button">Salvar aprendente</button>
	</div>
</form>
