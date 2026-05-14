<script lang="ts">
	import type { Learner, Visit, VisitStatus } from '$lib/modules/learners';
	import { formatCurrencyFromCents, formatDateForNBR, formatMonth } from '$lib/shared/formatters';

	let { learners } = $props<{
		learners: Learner[];
	}>();

	type FinanceVisit = {
		learner: Learner;
		visit: Visit;
		amountCents: number;
		monthKey: string;
		monthDate: Date;
	};

	type MonthBucket = {
		key: string;
		date: Date;
		label: string;
		revenueCents: number;
		completedCents: number;
		scheduledCents: number;
		missedCents: number;
		sessions: number;
		completed: number;
		scheduled: number;
		missed: number;
	};

	type LearnerRevenue = {
		id: string;
		name: string;
		revenueCents: number;
		sessions: number;
	};

	const today = new Date();
	const todayValue = formatDateForNBR(today);
	const currentMonthKey = getMonthKey(today);
	const currentMonthLabel = formatMonth(getMonthDateFromKey(currentMonthKey));

	const paidVisits = $derived(buildPaidVisits(learners));
	const monthBuckets = $derived(buildMonthBuckets(paidVisits));
	const currentMonth = $derived(
		monthBuckets.find((bucket: MonthBucket) => bucket.key === currentMonthKey) ??
			createMonthBucket(currentMonthKey)
	);
	const visibleMonths = $derived(monthBuckets.slice(0, 8));
	const maxMonthRevenue = $derived(
		Math.max(1, ...visibleMonths.map((bucket: MonthBucket) => bucket.revenueCents))
	);
	const learnerRevenue = $derived(buildLearnerRevenue(paidVisits));
	const maxLearnerRevenue = $derived(
		Math.max(1, ...learnerRevenue.map((bucket: LearnerRevenue) => bucket.revenueCents))
	);
	const nextBillings = $derived(
		paidVisits
			.filter((item: FinanceVisit) => isBillableStatus(item.visit.status) && item.visit.date >= todayValue)
			.slice(0, 8)
	);
	const paidLearnerCount = $derived(
		learners.filter((learner: Learner) => learner.sessionPriceCents > 0).length
	);
	const activeLearnersWithoutPrice = $derived(
		learners.filter(
			(learner: Learner) => learner.status === 'active' && learner.sessionPriceCents <= 0
		).length
	);
	const generalValueTotal = $derived(
		learners.reduce((sum: number, learner: Learner) => sum + learner.generalValueCents, 0)
	);
	const averageTicketCents = $derived(
		currentMonth.sessions > 0 ? Math.round(currentMonth.revenueCents / currentMonth.sessions) : 0
	);
	const kpis = $derived([
		{
			label: 'Faturamento do mes',
			value: formatCurrencyFromCents(currentMonth.revenueCents),
			detail: `${currentMonth.sessions} sessoes pagas previstas`,
			tone: 'purple'
		},
		{
			label: 'Realizado',
			value: formatCurrencyFromCents(currentMonth.completedCents),
			detail: `${currentMonth.completed} sessoes concluidas`,
			tone: 'green'
		},
		{
			label: 'A receber',
			value: formatCurrencyFromCents(currentMonth.scheduledCents),
			detail: `${currentMonth.scheduled} sessoes agendadas`,
			tone: 'blue'
		},
		{
			label: 'Valores gerais',
			value: formatCurrencyFromCents(generalValueTotal),
			detail: `${paidLearnerCount} aprendentes com valor por sessao`,
			tone: 'amber'
		}
	]);
	const statusCards = $derived([
		{
			label: 'Concluidas',
			value: currentMonth.completed,
			amountCents: currentMonth.completedCents,
			tone: 'green'
		},
		{
			label: 'Agendadas',
			value: currentMonth.scheduled,
			amountCents: currentMonth.scheduledCents,
			tone: 'blue'
		},
		{
			label: 'Perdidas',
			value: currentMonth.missed,
			amountCents: currentMonth.missedCents,
			tone: 'amber'
		}
	]);

	function buildPaidVisits(sourceLearners: Learner[]) {
		const items: FinanceVisit[] = [];

		for (const learner of sourceLearners) {
			if (learner.sessionPriceCents <= 0) continue;

			for (const visit of learner.visits) {
				const visitDate = parseVisitDate(visit.date);
				if (Number.isNaN(visitDate.getTime())) continue;

				items.push({
					learner,
					visit,
					amountCents: learner.sessionPriceCents,
					monthKey: getMonthKey(visitDate),
					monthDate: getMonthDateFromKey(getMonthKey(visitDate))
				});
			}
		}

		return items.sort((left, right) => {
			const dateOrder = left.visit.date.localeCompare(right.visit.date);
			if (dateOrder !== 0) return dateOrder;

			return left.visit.startTime.localeCompare(right.visit.startTime);
		});
	}

	function buildMonthBuckets(items: FinanceVisit[]) {
		const buckets = new Map<string, MonthBucket>();
		const currentMonthDate = getMonthDateFromKey(currentMonthKey);

		for (let index = 0; index < 6; index += 1) {
			const key = getMonthKey(addMonths(currentMonthDate, index));
			buckets.set(key, createMonthBucket(key));
		}

		for (const item of items) {
			if (item.monthKey < currentMonthKey) continue;

			const bucket = buckets.get(item.monthKey) ?? createMonthBucket(item.monthKey);
			buckets.set(item.monthKey, bucket);

			if (item.visit.status === 'missed') {
				bucket.missed += 1;
				bucket.missedCents += item.amountCents;
				continue;
			}

			bucket.sessions += 1;
			bucket.revenueCents += item.amountCents;

			if (item.visit.status === 'completed') {
				bucket.completed += 1;
				bucket.completedCents += item.amountCents;
			} else {
				bucket.scheduled += 1;
				bucket.scheduledCents += item.amountCents;
			}
		}

		return Array.from(buckets.values()).sort((left, right) => left.key.localeCompare(right.key));
	}

	function buildLearnerRevenue(items: FinanceVisit[]) {
		const buckets = new Map<string, LearnerRevenue>();

		for (const item of items) {
			if (item.monthKey !== currentMonthKey || !isBillableStatus(item.visit.status)) continue;

			const bucket =
				buckets.get(item.learner.id) ??
				({
					id: item.learner.id,
					name: item.learner.name,
					revenueCents: 0,
					sessions: 0
				} satisfies LearnerRevenue);

			bucket.revenueCents += item.amountCents;
			bucket.sessions += 1;
			buckets.set(item.learner.id, bucket);
		}

		return Array.from(buckets.values()).sort((left, right) => right.revenueCents - left.revenueCents);
	}

	function createMonthBucket(key: string): MonthBucket {
		const date = getMonthDateFromKey(key);

		return {
			key,
			date,
			label: formatMonth(date),
			revenueCents: 0,
			completedCents: 0,
			scheduledCents: 0,
			missedCents: 0,
			sessions: 0,
			completed: 0,
			scheduled: 0,
			missed: 0
		};
	}

	function parseVisitDate(value: string) {
		return new Date(`${value}T12:00:00`);
	}

	function getMonthKey(value: Date) {
		return `${value.getFullYear()}-${String(value.getMonth() + 1).padStart(2, '0')}`;
	}

	function getMonthDateFromKey(key: string) {
		return new Date(`${key}-01T12:00:00`);
	}

	function addMonths(value: Date, amount: number) {
		const next = new Date(value);
		next.setMonth(next.getMonth() + amount);
		return next;
	}

	function getPercent(value: number, maxValue: number) {
		if (maxValue <= 0) return 0;
		return Math.round((value / maxValue) * 100);
	}

	function isBillableStatus(status: VisitStatus) {
		return status !== 'missed';
	}

	function getVisitStatusLabel(status: VisitStatus) {
		const labels: Record<VisitStatus, string> = {
			scheduled: 'Agendada',
			completed: 'Concluida',
			missed: 'Perdida'
		};

		return labels[status];
	}
</script>

<section class="finance-workspace">
	<div class="finance-head">
		<div>
			<span class="section-kicker">Financeiro</span>
			<h1>Visao financeira</h1>
			<p>
				Projecao baseada nas sessoes com valor por sessao cadastrado em cada aprendente.
			</p>
		</div>

		<div class="finance-period-pill">
			<span>Mes atual</span>
			<strong>{currentMonthLabel}</strong>
		</div>
	</div>

	<div class="finance-kpi-grid">
		{#each kpis as metric}
			<article class={`finance-kpi ${metric.tone}`}>
				<span>{metric.label}</span>
				<strong>{metric.value}</strong>
				<p>{metric.detail}</p>
			</article>
		{/each}
	</div>

	<div class="finance-dashboard-grid">
		<article class="finance-panel finance-panel-wide">
			<header>
				<div>
					<h2>Receita prevista por mes</h2>
					<p>Mes atual e proximos meses com sessoes pagas.</p>
				</div>
			</header>

			<div class="finance-chart-list">
				{#each visibleMonths as month}
					<div class="finance-bar-row">
						<div class="finance-bar-copy">
							<strong>{month.label}</strong>
							<span>{month.sessions} sessoes faturaveis</span>
						</div>
						<div class="finance-bar-track" aria-hidden="true">
							<span
								class="finance-bar-fill"
								style={`width: ${getPercent(month.revenueCents, maxMonthRevenue)}%`}
							></span>
						</div>
						<strong>{formatCurrencyFromCents(month.revenueCents)}</strong>
					</div>
				{/each}
			</div>
		</article>

		<article class="finance-panel">
			<header>
				<div>
					<h2>Status do mes</h2>
					<p>{currentMonth.sessions + currentMonth.missed} sessoes com valor cadastrado.</p>
				</div>
			</header>

			<div class="finance-status-grid">
				{#each statusCards as status}
					<div class={`finance-status-card ${status.tone}`}>
						<span>{status.label}</span>
						<strong>{status.value}</strong>
						<small>{formatCurrencyFromCents(status.amountCents)}</small>
					</div>
				{/each}
			</div>

			<div class="finance-ratio-card">
				<span>Ticket medio do mes</span>
				<strong>{formatCurrencyFromCents(averageTicketCents)}</strong>
				<p>
					{activeLearnersWithoutPrice > 0
						? `${activeLearnersWithoutPrice} aprendentes ativos ainda sem valor por sessao.`
						: 'Todos os aprendentes ativos com preco entram na projecao.'}
				</p>
			</div>
		</article>
	</div>

	<div class="finance-bottom-grid">
		<article class="finance-panel">
			<header>
				<div>
					<h2>Receita por aprendente</h2>
					<p>Distribuicao prevista no mes atual.</p>
				</div>
			</header>

			<div class="finance-learner-chart">
				{#each learnerRevenue as item}
					<div class="finance-learner-row">
						<div>
							<strong>{item.name}</strong>
							<span>{item.sessions} sessoes</span>
						</div>
						<div class="finance-learner-track" aria-hidden="true">
							<span
								style={`width: ${getPercent(item.revenueCents, maxLearnerRevenue)}%`}
							></span>
						</div>
						<strong>{formatCurrencyFromCents(item.revenueCents)}</strong>
					</div>
				{:else}
					<p class="finance-empty">Nenhum faturamento previsto para o mes atual.</p>
				{/each}
			</div>
		</article>

		<article class="finance-panel">
			<header>
				<div>
					<h2>Proximos faturamentos</h2>
					<p>Sessoes futuras com valor por sessao.</p>
				</div>
			</header>

			<div class="finance-billing-list">
				{#each nextBillings as item}
					<div class="finance-billing-row">
						<div>
							<strong>{item.learner.name}</strong>
							<span>{formatDateForNBR(item.visit.date)} - {item.visit.startTime}</span>
						</div>
						<span class={`finance-mini-status ${item.visit.status}`}>
							{getVisitStatusLabel(item.visit.status)}
						</span>
						<strong>{formatCurrencyFromCents(item.amountCents)}</strong>
					</div>
				{:else}
					<p class="finance-empty">Nao ha sessoes futuras pagas cadastradas.</p>
				{/each}
			</div>
		</article>
	</div>
</section>
