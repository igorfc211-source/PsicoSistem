<script lang="ts">
	import { fly } from 'svelte/transition';
	import type {
		CommunicationFamily,
		CommunicationContact,
		CommunicationResponsible,
		CommunicationStage,
		ContactChannelType,
		GuardianOption,
		Learner,
		NewCommunicationContactInput,
		NewCommunicationFamilyInput,
		NewCommunicationResponsibleInput
	} from '$lib/modules/learners';
	import {
		COMMUNICATION_STAGES,
		CONTACT_CHANNEL_TYPES,
		RELATIONSHIP_OPTIONS,
		buildGuardianOptionsFromLearners,
		getContactTypeLabel,
		getCommunicationFamilyResponsibleKeys,
		getFamilyLearners,
		normalizeInstagramHandle,
		normalizePhone
	} from '$lib/modules/learners';

	let {
		learners,
		families,
		selectedFamilyId,
		searchTerm,
		onCreateFamily,
		onUpdateFamily,
		onDeleteFamily,
		onAddResponsible,
		onRemoveResponsible,
		onAddContact,
		onRemoveContact,
		onSelectFamily,
		onCloseFamily
	} = $props<{
		learners: Learner[];
		families: CommunicationFamily[];
		selectedFamilyId: string | null;
		searchTerm: string;
		onCreateFamily: (input: NewCommunicationFamilyInput) => boolean;
		onUpdateFamily: (familyId: string, patch: Partial<CommunicationFamily>) => void;
		onDeleteFamily: (familyId: string) => void;
		onAddResponsible: (familyId: string, input: NewCommunicationResponsibleInput) => boolean;
		onRemoveResponsible: (familyId: string, responsibleId: string) => void;
		onAddContact: (familyId: string, input: NewCommunicationContactInput) => boolean;
		onRemoveContact: (familyId: string, contactId: string) => void;
		onSelectFamily: (familyId: string) => void;
		onCloseFamily: () => void;
	}>();

	let activeStage = $state<CommunicationStage | 'all'>('all');
	let showFamilyForm = $state(false);
	let composerFamilyId = $state<string | null>(null);
	let composerType = $state<'responsible' | 'contact' | null>(null);
	let familyInput = $state<NewCommunicationFamilyInput>(createEmptyFamilyInput());
	let responsibleInput = $state<NewCommunicationResponsibleInput>(createEmptyResponsibleInput());
	let contactInput = $state<NewCommunicationContactInput>(createEmptyContactInput());
	let learnerToAddId = $state('');
	let selectedFamilyGuardianKey = $state('');
	let selectedResponsibleGuardianKey = $state('');

	const filteredFamilies = $derived(
		families
			.filter((family: CommunicationFamily) => {
				if (activeStage !== 'all' && family.stage !== activeStage) return false;
				const query = searchTerm.trim().toLowerCase();
				if (!query) return true;

				const linkedLearners = getFamilyLearners(family, learners);
				const searchableText = [
					family.familyName,
					...family.responsibles.flatMap((responsible: CommunicationResponsible) => [
						responsible.name,
						responsible.relationship,
						responsible.phone
					]),
					...family.contacts.flatMap((contact: CommunicationContact) => [
						contact.label,
						contact.value,
						contact.notes
					]),
					...linkedLearners.map((learner) => learner.name)
				]
					.join(' ')
					.toLowerCase();

				return searchableText.includes(query);
			})
			.sort((left: CommunicationFamily, right: CommunicationFamily) =>
				right.updatedAt.localeCompare(left.updatedAt)
			)
	);
	const linkedLearnerCount = $derived(
		new Set(families.flatMap((family: CommunicationFamily) => family.learnerIds)).size
	);
	const availableLearnersForFamilyInput = $derived(
		learners.filter((learner: Learner) => !familyInput.learnerIds.includes(learner.id))
	);
	const selectedLearnersForFamilyInput = $derived(
		learners.filter((learner: Learner) => familyInput.learnerIds.includes(learner.id))
	);
	const guardianOptions = $derived(buildGuardianOptionsFromLearners(learners, families));
	const usedResponsibleKeys = $derived(
		new Set(families.flatMap((family: CommunicationFamily) => getCommunicationFamilyResponsibleKeys(family)))
	);
	const availableFamilyGuardianOptions = $derived(
		guardianOptions.filter((option: GuardianOption) => !usedResponsibleKeys.has(option.key))
	);

	function submitFamily(event: Event) {
		event.preventDefault();
		if (!onCreateFamily(familyInput)) return;

		resetFamilyForm();
		showFamilyForm = false;
	}

	function submitResponsible(event: Event, familyId: string) {
		event.preventDefault();
		if (!onAddResponsible(familyId, responsibleInput)) return;

		responsibleInput = createEmptyResponsibleInput();
		selectedResponsibleGuardianKey = '';
		closeComposer();
	}

	function submitContact(event: Event, familyId: string) {
		event.preventDefault();
		if (!onAddContact(familyId, contactInput)) return;

		contactInput = createEmptyContactInput();
		closeComposer();
	}

	function openResponsibleComposer(familyId: string) {
		onSelectFamily(familyId);
		composerFamilyId = familyId;
		composerType = 'responsible';
		responsibleInput = createEmptyResponsibleInput();
		selectedResponsibleGuardianKey = '';
	}

	function openContactComposer(familyId: string) {
		onSelectFamily(familyId);
		composerFamilyId = familyId;
		composerType = 'contact';
		contactInput = createEmptyContactInput();
	}

	function closeComposer() {
		composerFamilyId = null;
		composerType = null;
		selectedResponsibleGuardianKey = '';
	}

	function closeFamilyCard() {
		closeComposer();
		onCloseFamily();
	}

	function toggleFamilyForm() {
		showFamilyForm = !showFamilyForm;
		if (showFamilyForm) resetFamilyForm();
	}

	function resetFamilyForm() {
		familyInput = createEmptyFamilyInput();
		learnerToAddId = '';
		selectedFamilyGuardianKey = '';
	}

	function addSelectedLearnerToFamilyInput() {
		if (!learnerToAddId || familyInput.learnerIds.includes(learnerToAddId)) return;

		familyInput = {
			...familyInput,
			learnerIds: [...familyInput.learnerIds, learnerToAddId]
		};
		learnerToAddId = '';
	}

	function handleLearnerSelection(event: Event) {
		learnerToAddId = (event.currentTarget as HTMLSelectElement).value;
		addSelectedLearnerToFamilyInput();
	}

	function removeLearnerFromFamilyInput(learnerId: string) {
		familyInput = {
			...familyInput,
			learnerIds: familyInput.learnerIds.filter((id) => id !== learnerId)
		};
	}

	function handleFamilyGuardianSelection(event: Event) {
		const key = (event.currentTarget as HTMLSelectElement).value;
		selectedFamilyGuardianKey = key;
		const option = guardianOptions.find((item: GuardianOption) => item.key === key);
		if (!option) {
			familyInput = {
				...familyInput,
				responsibleName: '',
				relationship: '',
				learnerIds: []
			};
			return;
		}

		familyInput = {
			...familyInput,
			familyName: familyInput.familyName || inferFamilyName(option.name),
			responsibleName: option.name,
			responsiblePhone: option.phone,
			relationship: option.relationship,
			learnerIds: Array.from(new Set([...familyInput.learnerIds, ...option.learnerIds]))
		};
	}

	function handleResponsibleGuardianSelection(event: Event) {
		const key = (event.currentTarget as HTMLSelectElement).value;
		selectedResponsibleGuardianKey = key;
		const option = guardianOptions.find((item: GuardianOption) => item.key === key);
		if (!option) {
			responsibleInput = createEmptyResponsibleInput();
			return;
		}

		responsibleInput = {
			name: option.name,
			relationship: option.relationship,
			phone: option.phone,
			learnerIds: option.learnerIds
		};
	}

	function updateFamilyResponsiblePhone(event: Event) {
		const phone = normalizePhone((event.currentTarget as HTMLInputElement).value).slice(0, 11);
		familyInput = { ...familyInput, responsiblePhone: phone };
	}

	function updateResponsiblePhone(event: Event) {
		const phone = normalizePhone((event.currentTarget as HTMLInputElement).value).slice(0, 11);
		responsibleInput = { ...responsibleInput, phone };
	}

	function updateContactType(event: Event) {
		const type = (event.currentTarget as HTMLSelectElement).value as ContactChannelType;
		contactInput = {
			type,
			label: '',
			value: '',
			notes: ''
		};
	}

	function updateContactPhone(event: Event) {
		const value = normalizePhone((event.currentTarget as HTMLInputElement).value).slice(0, 11);
		contactInput = { ...contactInput, value };
	}

	function updateContactInstagram(event: Event) {
		const value = normalizeInstagramHandle((event.currentTarget as HTMLInputElement).value);
		contactInput = { ...contactInput, value };
	}

	function updateFamilyStage(familyId: string, event: Event) {
		const stage = (event.currentTarget as HTMLSelectElement).value as CommunicationStage;
		onUpdateFamily(familyId, { stage });
	}

	function createEmptyFamilyInput(): NewCommunicationFamilyInput {
		return {
			familyName: '',
			responsibleName: '',
			responsiblePhone: '',
			relationship: '',
			learnerIds: []
		};
	}

	function createEmptyResponsibleInput(): NewCommunicationResponsibleInput {
		return {
			name: '',
			relationship: '',
			phone: '',
			learnerIds: []
		};
	}

	function createEmptyContactInput(): NewCommunicationContactInput {
		return {
			type: 'whatsapp',
			label: '',
			value: '',
			notes: ''
		};
	}

	function getStageLabel(stage: CommunicationStage) {
		return COMMUNICATION_STAGES.find((item) => item.value === stage)?.label ?? 'Novo';
	}

	function formatPhone(phone: string) {
		const digits = phone.replace(/\D/g, '');
		if (digits.length === 11) {
			return `(${digits.slice(0, 2)}) ${digits.slice(2, 7)}-${digits.slice(7)}`;
		}
		if (digits.length === 10) {
			return `(${digits.slice(0, 2)}) ${digits.slice(2, 6)}-${digits.slice(6)}`;
		}
		return phone || 'Sem numero';
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

	function formatContactValue(contact: CommunicationContact) {
		if (contact.type === 'phone' || contact.type === 'whatsapp') return formatPhone(contact.value);
		return contact.value;
	}

	function getContactValueLabel(type: ContactChannelType) {
		if (type === 'phone' || type === 'whatsapp') return 'Telefone';
		if (type === 'instagram') return '@ da conta';
		if (type === 'email') return 'E-mail';
		return 'Contato';
	}

	function getContactValuePlaceholder(type: ContactChannelType) {
		if (type === 'phone' || type === 'whatsapp') return '(11) 99999-9999';
		if (type === 'instagram') return '@familiasilva';
		if (type === 'email') return 'familia@email.com';
		return 'Site, usuario ou outro contato';
	}

	function getContactLabelPlaceholder(type: ContactChannelType) {
		if (type === 'instagram') return 'Nome do perfil';
		if (type === 'email') return 'Nome do contato';
		if (type === 'phone' || type === 'whatsapp') return 'Nome de contato';
		return 'Nome do contato';
	}

	function getResponsibleSummary(family: CommunicationFamily) {
		const names = family.responsibles
			.map((responsible: CommunicationResponsible) => responsible.name.trim())
			.filter(Boolean);

		return formatCompactNameSummary(names, 'Sem responsavel');
	}

	function getLearnerSummary(linkedLearners: Learner[]) {
		const names = linkedLearners.map((learner) => learner.name.trim()).filter(Boolean);

		return formatCompactNameSummary(names, 'Nenhum aprendente vinculado');
	}

	function formatCompactNameSummary(names: string[], fallback: string) {
		if (!names.length) return fallback;
		if (names.length === 1) return names[0];

		return `${names[0]} +${names.length - 1}`;
	}

	function getAvailableResponsibleOptions(family: CommunicationFamily) {
		const blockedKeys = new Set<string>();
		for (const item of families) {
			if (item.id === family.id) continue;
			for (const key of getCommunicationFamilyResponsibleKeys(item)) blockedKeys.add(key);
		}

		const currentFamilyKeys = new Set(getCommunicationFamilyResponsibleKeys(family));
		return guardianOptions.filter(
			(option: GuardianOption) => !blockedKeys.has(option.key) && !currentFamilyKeys.has(option.key)
		);
	}

	function inferFamilyName(responsibleName: string) {
		const tokens = responsibleName.split(/\s+/).filter(Boolean);
		const surname = tokens.at(-1) ?? 'Contato';

		return `Familia ${surname}`;
	}

	function getContactHref(type: ContactChannelType, value: string) {
		const cleanValue = value.trim();
		if (!cleanValue) return '';

		if (type === 'email') return `mailto:${cleanValue}`;
		if (type === 'phone') return `tel:${cleanValue.replace(/\D/g, '')}`;
		if (type === 'whatsapp') return `https://wa.me/55${cleanValue.replace(/\D/g, '')}`;
		if (type === 'instagram') {
			const profile = cleanValue.replace(/^@/, '').replace(/^https?:\/\/(www\.)?instagram\.com\//, '');
			return `https://instagram.com/${profile}`;
		}

		return cleanValue.startsWith('http') ? cleanValue : '';
	}
</script>

<section class="communications-workspace">
	<div class="communication-head">
		<div>
			<h1>Comunicacoes</h1>
			<p>{families.length} familias em relacionamento</p>
		</div>

		<div class="title-actions">
			<button type="button" class="primary-button" onclick={toggleFamilyForm}>
				+ Familia
			</button>
		</div>
	</div>

	{#if showFamilyForm}
		<form
			class="communication-form family-creator"
			onsubmit={submitFamily}
			transition:fly={{ y: -14, duration: 180 }}
		>
			<div class="form-grid compact">
				<label>
					<span>Nome do card</span>
					<input bind:value={familyInput.familyName} placeholder="Familia Silva" required />
				</label>
				<label>
					<span>Responsavel disponivel</span>
					<select
						value={selectedFamilyGuardianKey}
						onchange={handleFamilyGuardianSelection}
						required
						disabled={!availableFamilyGuardianOptions.length}
					>
						<option value="">
							{availableFamilyGuardianOptions.length
								? 'Selecione um responsavel'
								: 'Todos os responsaveis ja possuem card'}
						</option>
						{#each availableFamilyGuardianOptions as option}
							<option value={option.key}>
								{option.name} - {option.learnerIds.length} aprendente{option.learnerIds.length === 1 ? '' : 's'}
							</option>
						{/each}
					</select>
				</label>
				<label>
					<span>Responsavel principal</span>
					<input bind:value={familyInput.responsibleName} placeholder="Selecione acima" readonly required />
				</label>
				<label>
					<span>Numero</span>
					<input
						value={formatPhoneInput(familyInput.responsiblePhone)}
						inputmode="numeric"
						pattern={'\\([0-9]{2}\\) [0-9]{4,5}-[0-9]{4}'}
						maxlength="15"
						required
						placeholder="(11) 99999-9999"
						title="Informe um numero no formato (xx) xxxxx-xxxx."
						oninput={updateFamilyResponsiblePhone}
					/>
				</label>
				<label>
					<span>Parentesco</span>
					<select bind:value={familyInput.relationship}>
						<option value="">Selecione</option>
						{#each RELATIONSHIP_OPTIONS as option}
							<option value={option}>{option}</option>
						{/each}
					</select>
				</label>
			</div>

			<div class="learner-picker">
				<span>Aprendentes vinculados</span>
				<div class="learner-select-row">
					<select value={learnerToAddId} aria-label="Selecionar aprendente" onchange={handleLearnerSelection}>
						<option value="">Selecione um aprendente</option>
						{#each availableLearnersForFamilyInput as learner}
							<option value={learner.id}>{learner.name}</option>
						{/each}
					</select>
				</div>
				<div class="selected-learner-list">
					{#if selectedLearnersForFamilyInput.length}
						{#each selectedLearnersForFamilyInput as learner}
							<span>
								{learner.name}
								<button
									type="button"
									aria-label={`Remover ${learner.name}`}
									onclick={() => removeLearnerFromFamilyInput(learner.id)}
								>
									x
								</button>
							</span>
						{/each}
					{:else}
						<p>Nenhum aprendente vinculado a este card.</p>
					{/if}
				</div>
			</div>

			<div class="form-actions">
				<button
					type="button"
					class="secondary-button"
					onclick={() => {
						resetFamilyForm();
						showFamilyForm = false;
					}}
				>
					Cancelar
				</button>
				<button type="submit" class="primary-button" disabled={!availableFamilyGuardianOptions.length}>
					Criar card
				</button>
			</div>
		</form>
	{/if}

	<div class="crm-stage-tabs" aria-label="Filtrar funil de comunicacao">
		<button type="button" class:active={activeStage === 'all'} onclick={() => (activeStage = 'all')}>
			Todos
		</button>
		{#each COMMUNICATION_STAGES as stage}
			<button
				type="button"
				class:active={activeStage === stage.value}
				onclick={() => (activeStage = stage.value)}
			>
				{stage.label}
			</button>
		{/each}
	</div>

	{#if filteredFamilies.length === 0}
		<div class="communication-empty">
			<h2>Nenhuma familia encontrada</h2>
			<p>Crie um card de comunicacao ou cadastre um aprendente com responsavel.</p>
		</div>
	{:else}
		<div class="family-card-grid">
			{#each filteredFamilies as family (family.id)}
				{@const linkedLearners = getFamilyLearners(family, learners)}
				{@const isSelected = selectedFamilyId === family.id}
				<article
					class="family-card"
					class:selected={isSelected}
				>
					<button
						type="button"
						class="family-card-summary"
						class:expanded={isSelected}
						aria-expanded={isSelected}
						onclick={() => (isSelected ? closeFamilyCard() : onSelectFamily(family.id))}
					>
						<span>Familia</span>
						<span class="family-status-badge">{getStageLabel(family.stage)}</span>
						<span class="family-summary-title">{family.familyName}</span>
						<span class="family-summary-line">
							<span>Responsavel</span>
							<strong>{getResponsibleSummary(family)}</strong>
						</span>
						<span class="family-summary-line">
							<span>Aprendente</span>
							<strong>{getLearnerSummary(linkedLearners)}</strong>
						</span>
					</button>

					{#if isSelected}
						<div class="family-card-details" transition:fly={{ y: -10, duration: 180 }}>
							<header class="family-card-toolbar">
								<div>
									<span>Estagio</span>
									<strong>{getStageLabel(family.stage)}</strong>
								</div>
								<div class="family-card-actions">
									<select
										value={family.stage}
										onchange={(event) => updateFamilyStage(family.id, event)}
										aria-label="Estagio da familia"
									>
										{#each COMMUNICATION_STAGES as stage}
											<option value={stage.value}>{stage.label}</option>
										{/each}
									</select>
									
								</div>
							</header>

					<section class="family-block">
						<div class="family-block-head">
							<h3>Responsaveis</h3>
							<button
								type="button"
								class="mini-action"
								onclick={() => openResponsibleComposer(family.id)}
							>
								+ Responsavel
							</button>
						</div>

						<div class="contact-list">
							{#each family.responsibles as responsible}
								<div class="contact-row">
									<div>
										<strong>{responsible.name}</strong>
										<span>{responsible.relationship || 'Responsavel'}</span>
									</div>
									<a href={`tel:${responsible.phone}`}>{formatPhone(responsible.phone)}</a>
									<button
										type="button"
										aria-label="Remover responsavel"
										onclick={() => onRemoveResponsible(family.id, responsible.id)}
									>
										x
									</button>
								</div>
							{/each}
						</div>

						{#if composerFamilyId === family.id && composerType === 'responsible'}
							{@const responsibleOptions = getAvailableResponsibleOptions(family)}
							<form
								class="inline-composer"
								onsubmit={(event) => submitResponsible(event, family.id)}
							>
								<div class="form-grid compact">
									<label>
										<span>Responsavel disponivel</span>
										<select
											value={selectedResponsibleGuardianKey}
											onchange={handleResponsibleGuardianSelection}
											required
											disabled={!responsibleOptions.length || family.responsibles.length >= 2}
										>
											<option value="">
												{family.responsibles.length >= 2
													? 'Limite de dois responsaveis'
													: responsibleOptions.length
														? 'Selecione um responsavel'
														: 'Nenhum responsavel disponivel'}
											</option>
											{#each responsibleOptions as option}
												<option value={option.key}>
													{option.name} - {option.learnerIds.length} aprendente{option.learnerIds.length === 1
														? ''
														: 's'}
												</option>
											{/each}
										</select>
									</label>
									<label>
										<span>Nome</span>
										<input bind:value={responsibleInput.name} placeholder="Selecione acima" readonly required />
									</label>
									<label>
										<span>Numero</span>
										<input
											value={formatPhoneInput(responsibleInput.phone)}
											inputmode="numeric"
											pattern={'\\([0-9]{2}\\) [0-9]{4,5}-[0-9]{4}'}
											maxlength="15"
											required
											placeholder="(11) 99999-9999"
											title="Informe um numero no formato (xx) xxxxx-xxxx."
											oninput={updateResponsiblePhone}
										/>
									</label>
									<label>
										<span>Parentesco</span>
										<select bind:value={responsibleInput.relationship}>
											<option value="">Selecione</option>
											{#each RELATIONSHIP_OPTIONS as option}
												<option value={option}>{option}</option>
											{/each}
										</select>
									</label>
								</div>
								<div class="form-actions">
									<button type="button" class="secondary-button" onclick={closeComposer}>Cancelar</button>
									<button
										type="submit"
										class="primary-button"
										disabled={!responsibleOptions.length || family.responsibles.length >= 2}
									>
										Adicionar
									</button>
								</div>
							</form>
						{/if}
					</section>

					<section class="family-block">
						<div class="family-block-head">
							<h3>Canais</h3>
							<button
								type="button"
								class="mini-action"
								onclick={() => openContactComposer(family.id)}
							>
								+ Contato
							</button>
						</div>

						<div class="contact-list">
							{#each family.contacts as contact}
								<div class="contact-row">
									<div>
										<strong>{contact.label || getContactTypeLabel(contact.type)}</strong>
										<span>{getContactTypeLabel(contact.type)}</span>
									</div>
									{#if getContactHref(contact.type, contact.value)}
										<a href={getContactHref(contact.type, contact.value)} target="_blank" rel="noreferrer">
											{formatContactValue(contact)}
										</a>
									{:else}
										<span>{formatContactValue(contact)}</span>
									{/if}
									<button
										type="button"
										aria-label="Remover contato"
										onclick={() => onRemoveContact(family.id, contact.id)}
									>
										x
									</button>
								</div>
							{/each}
						</div>

						{#if composerFamilyId === family.id && composerType === 'contact'}
							<form
								class="inline-composer"
								onsubmit={(event) => submitContact(event, family.id)}
							>
								<div class="form-grid compact">
									<label>
										<span>Tipo</span>
										<select value={contactInput.type} onchange={updateContactType}>
											{#each CONTACT_CHANNEL_TYPES as type}
												<option value={type.value}>{type.label}</option>
											{/each}
										</select>
									</label>
									<label>
										<span>{getContactValueLabel(contactInput.type)}</span>
										{#if contactInput.type === 'phone' || contactInput.type === 'whatsapp'}
											<input
												value={formatPhoneInput(contactInput.value)}
												inputmode="numeric"
												pattern={'\\([0-9]{2}\\) [0-9]{4,5}-[0-9]{4}'}
												maxlength="15"
												required
												placeholder={getContactValuePlaceholder(contactInput.type)}
												title="Informe um numero no formato (xx) xxxxx-xxxx."
												oninput={updateContactPhone}
											/>
										{:else if contactInput.type === 'instagram'}
											<input
												value={contactInput.value}
												required
												placeholder={getContactValuePlaceholder(contactInput.type)}
												title="Informe o @ da conta."
												oninput={updateContactInstagram}
											/>
										{:else if contactInput.type === 'email'}
											<input
												type="email"
												bind:value={contactInput.value}
												required
												placeholder={getContactValuePlaceholder(contactInput.type)}
											/>
										{:else}
											<input
												bind:value={contactInput.value}
												required
												placeholder={getContactValuePlaceholder(contactInput.type)}
											/>
										{/if}
									</label>
									<label>
										<span>Nome do contato</span>
										<input
											bind:value={contactInput.label}
											required={contactInput.type !== 'instagram'}
											placeholder={getContactLabelPlaceholder(contactInput.type)}
										/>
									</label>
								</div>
								<label>
									<span>Observacao</span>
									<input bind:value={contactInput.notes} placeholder="Melhor canal para lembretes" />
								</label>
								<div class="form-actions">
									<button type="button" class="secondary-button" onclick={closeComposer}>Cancelar</button>
									<button type="submit" class="primary-button">Adicionar</button>
								</div>
							</form>
						{/if}
					</section>

					<footer>
						<span>{linkedLearners.length} aprendente{linkedLearners.length === 1 ? '' : 's'}</span>
						<button
							type="button"
							class="danger-link"
							onclick={() => onDeleteFamily(family.id)}
						>
							Remover card
						</button>
					</footer>
						</div>
					{/if}
				</article>
			{/each}
		</div>
	{/if}
</section>
