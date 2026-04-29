<script lang="ts">
	import { createDefaultLearnerInput, type NewLearnerInput } from '$lib/modules/learners';

	let {
		onCreate,
		onCancel
	} = $props<{
		onCreate: (input: NewLearnerInput) => boolean;
		onCancel: () => void;
	}>();

	let form = $state<NewLearnerInput>(createDefaultLearnerInput());

	// Envia uma copia do formulario para a pagina orquestradora validar e persistir.
	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		if (onCreate({ ...form })) {
			form = createDefaultLearnerInput();
		}
	}
</script>

<form class="add-form" onsubmit={handleSubmit}>
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

	<div class="form-actions">
		<button type="button" class="secondary-button" onclick={onCancel}>Cancelar</button>
		<button type="submit" class="primary-button">Salvar aprendente</button>
	</div>
</form>
