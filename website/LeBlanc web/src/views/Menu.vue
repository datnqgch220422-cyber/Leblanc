<script setup>
import { computed, onMounted, ref } from "vue";
import { getDrinks } from "@/api";
import { useThemeState } from "@/composables/useThemeState";
import { useCart } from "@/composables/useCart";
import americanoCamSa from "@/assets/americano-cam-sa.png";
import cafeSuaDa from "@/assets/cafe-sua-da.png";
import cafeDua from "@/assets/cafe-dua.png";
import latteYenMach from "@/assets/latte-sua-yen-mach.png";
import traDaoCamSa from "@/assets/tra-dao-cam-sa.png";
import traGungMatOng from "@/assets/tra-gung-mat-ong.png";
import nuocEpCamCaRot from "@/assets/nuoc-ep-cam-ca-rot.png";
import sinhToChuoiYenMach from "@/assets/sinh-to-chuoi-yen-mach.png";
import suaChuaChanhDay from "@/assets/sua-chua-chanh-day.png";
import matchaLatte from "@/assets/matcha-latte.png";
import caCaoQueNong from "@/assets/ca-cao-que-nong.png";
import traTaoQueMatOngNong from "@/assets/tra-tao-que-mat-ong-nong.png";
import ginTonicChanhBuoi from "@/assets/gin-tonic-chanh-buoi.png";
import mojitoBacHaCoDien from "@/assets/mojito-bac-ha-co-dien.png";
import rumCokeChanh from "@/assets/rum-coke-chanh.png";
import whiskyHighball from "@/assets/whisky-highball.png";
import houseRed from "@/assets/house-red.png";
import houseWhite from "@/assets/house-white.png";
import beer from "@/assets/beer.png";
import craftedBeer from "@/assets/crafted-beer.png";
import plumLiqueur from "@/assets/plum-liqueur.png";
import midnightCoffee from "@/assets/midnight-coffee.png";

// Morning: curated 12-item set with real photos.
const morningMenu = [
  {
    _id: "cafe-sua-da",
    name: "Vietnamese iced milk coffee",
    desc: "Robusta phin from Buon Ma Thuot, condensed milk, ice; bold and creamy in the Saigon style.",
    price: 29000,
    image: cafeSuaDa,
    caffeine: "~90 mg/cup",
    temp: "Cold (ice)",
    sweetness: 4,
    kind: "Phin coffee",
    isAlcoholic: false,
  },
  {
    _id: "americano-cam-sa",
    name: "Orange lemongrass americano",
    desc: "Da Lat Arabica espresso with fresh orange juice and smashed lemongrass; no milk; light citrus aroma.",
    price: 39000,
    image: americanoCamSa,
    caffeine: "~80 mg/cup",
    temp: "Cold (ice)",
    sweetness: 2,
    kind: "Espresso bar",
    isAlcoholic: false,
  },
  {
    _id: "latte-yen-mach",
    name: "Oat milk latte",
    desc: "Arabica espresso with unsweetened oat milk, fine foam; fits dairy-free guests.",
    price: 45000,
    image: latteYenMach,
    caffeine: "~75 mg/cup",
    temp: "Hot",
    sweetness: 3,
    kind: "Espresso bar",
    isAlcoholic: false,
  },
  {
    _id: "cafe-dua",
    name: "Coconut iced coffee",
    desc: "Arabica-Robusta espresso blend with light coconut cream and low-sugar milk, blended with ice.",
    price: 49000,
    image: cafeDua,
    caffeine: "~85 mg/cup",
    temp: "Cold (ice)",
    sweetness: 4,
    kind: "Signature coffee",
    isAlcoholic: false,
  },
  {
    _id: "tra-dao-cam-sa",
    name: "Peach orange lemongrass tea",
    desc: "Cold-brew Ceylon black tea with peach syrup, orange slices, lemongrass, peach chunks.",
    price: 39000,
    image: traDaoCamSa,
    caffeine: "-",
    temp: "Cold (ice)",
    sweetness: 3,
    kind: "Fruit tea",
    isAlcoholic: false,
  },
  {
    _id: "tra-gung-mat-ong",
    name: "Ginger honey tea",
    desc: "Light green tea, fresh ginger slices, forest honey; warming for morning.",
    price: 32000,
    image: traGungMatOng,
    caffeine: "-",
    temp: "Hot",
    sweetness: 2,
    kind: "Hot tea",
    isAlcoholic: false,
  },
  {
    _id: "nuoc-ep-cam-ca-rot",
    name: "Orange carrot juice",
    desc: "Fresh orange + carrot juice, no syrup; sweetness adjustable on request.",
    price: 39000,
    image: nuocEpCamCaRot,
    caffeine: "-",
    temp: "Cold (light ice)",
    sweetness: "2-3",
    kind: "Cold-pressed juice",
    isAlcoholic: false,
  },
  {
    _id: "sinh-to-chuoi-yen-mach",
    name: "Banana oat smoothie",
    desc: "Ripe banana, rolled oats, almond milk, touch of honey; filling and gym-friendly.",
    price: 45000,
    image: sinhToChuoiYenMach,
    caffeine: "-",
    temp: "Cold (blended)",
    sweetness: 3,
    kind: "Healthy smoothie",
    isAlcoholic: false,
  },
  {
    _id: "sua-chua-chanh-day",
    name: "Passionfruit yogurt frappe",
    desc: "Fermented yogurt with fresh passionfruit sauce, blended with ice; tangy-sweet.",
    price: 39000,
    image: suaChuaChanhDay,
    caffeine: "-",
    temp: "Cold (blended)",
    sweetness: "3-4",
    kind: "Yogurt",
    isAlcoholic: false,
  },
  {
    _id: "matcha-latte",
    name: "Matcha latte",
    desc: "Japanese matcha whisked with fresh milk, steamed; aromatic with a gentle bitterness.",
    price: 49000,
    image: matchaLatte,
    caffeine: "- (tea-based)",
    temp: "Hot",
    sweetness: 3,
    kind: "Tea latte",
    isAlcoholic: false,
  },
  {
    _id: "ca-cao-que-nong",
    name: "Hot cacao with cinnamon",
    desc: "Natural cacao powder, fresh milk (can add condensed milk), syrup or honey, plus cinnamon dust or a stick for warmth.",
    price: 45000,
    image: caCaoQueNong,
    caffeine: "-",
    temp: "Hot",
    sweetness: "3-4",
    kind: "Cacao drink",
    isAlcoholic: false,
  },
  {
    _id: "tra-tao-que-mat-ong-nong",
    name: "Warm apple cinnamon honey tea",
    desc: "Light black or green tea with fresh apple slices, cinnamon stick, honey; optionally a squeeze of lime; gentle apple-cinnamon aroma.",
    price: 39000,
    image: traTaoQueMatOngNong,
    caffeine: "-",
    temp: "Hot",
    sweetness: "2-3",
    kind: "Healthy hot tea",
    isAlcoholic: false,
  },
];

// Night: 10-item cocktail/boozy menu (English), with real photos.
// Signature drinks first: Midnight Coffee, Citrus Gin & Tonic, Classic Mint Mojito.
const nightSpecials = [
  {
    _id: "midnight-coffee",
    name: "Midnight Coffee",
    desc: "Coffee cocktail with rum/whisky and espresso; warm, boozy coffee sweetness.",
    price: 129000,
    image: midnightCoffee,
    caffeine: "~60 mg/cup",
    temp: "Cold",
    sweetness: "3-4",
    kind: "Signature coffee cocktail",
    isAlcoholic: true,
    baseSpirit: "Rum or whisky + espresso",
  },
  {
    _id: "gin-tonic-chanh-buoi",
    name: "Citrus Gin & Tonic",
    desc: "Gin & tonic with lemon and grapefruit peel; crisp, aromatic, lightly fizzy.",
    price: 99000,
    image: ginTonicChanhBuoi,
    caffeine: "-",
    temp: "Cold",
    sweetness: "Light",
    kind: "Signature cocktail",
    isAlcoholic: true,
    baseSpirit: "Gin + tonic",
  },
  {
    _id: "mojito-bac-ha-co-dien",
    name: "Classic Mint Mojito",
    desc: "White rum, fresh mint, lime, soda; cool and gently sweet.",
    price: 109000,
    image: mojitoBacHaCoDien,
    caffeine: "-",
    temp: "Cold",
    sweetness: "2-3",
    kind: "Signature cocktail",
    isAlcoholic: true,
    baseSpirit: "Rum + soda",
  },
  {
    _id: "rum-coke-chanh",
    name: "Rum & Coke with lime",
    desc: "White rum and cola with lime; familiar, bubbly, easy-drinking.",
    price: 89000,
    image: rumCokeChanh,
    caffeine: "-",
    temp: "Cold",
    sweetness: "3",
    kind: "Long drink",
    isAlcoholic: true,
    baseSpirit: "Rum + cola",
  },
  {
    _id: "whisky-highball",
    name: "Whisky Highball",
    desc: "Whisky with soda, subtle fizz; mellow malt aroma, not harsh.",
    price: 129000,
    image: whiskyHighball,
    caffeine: "-",
    temp: "Cold",
    sweetness: "2",
    kind: "Long drink",
    isAlcoholic: true,
    baseSpirit: "Whisky + soda",
  },
  {
    _id: "house-red",
    name: "House Red Wine (by the glass)",
    desc: "Fruity red wine, slightly tart; easy, dinner-friendly.",
    price: 119000,
    image: houseRed,
    caffeine: "-",
    temp: "Room temp",
    sweetness: "Light",
    kind: "Wine",
    isAlcoholic: true,
    baseSpirit: "Red wine",
  },
  {
    _id: "house-white",
    name: "House White Wine (by the glass)",
    desc: "Bright white wine with mild fruit (apple, pear, citrus).",
    price: 119000,
    image: houseWhite,
    caffeine: "-",
    temp: "Cold",
    sweetness: "Light",
    kind: "Wine",
    isAlcoholic: true,
    baseSpirit: "White wine",
  },
  {
    _id: "beer",
    name: "Classic lager beer",
    desc: "Cold lager, lightly bitter, crisp and easy.",
    price: 49000,
    image: beer,
    caffeine: "-",
    temp: "Cold",
    sweetness: "Low",
    kind: "Beer",
    isAlcoholic: true,
    baseSpirit: "Lager",
  },
  {
    _id: "crafted-beer",
    name: "Craft beer (1–2 rotating taps)",
    desc: "Pale ale / IPA / wheat; mildly bitter, hoppy or fruity depending on tap.",
    price: 89000,
    image: craftedBeer,
    caffeine: "-",
    temp: "Cold",
    sweetness: "Low",
    kind: "Craft beer",
    isAlcoholic: true,
    baseSpirit: "Beer",
  },
  {
    _id: "plum-liqueur",
    name: "Plum liqueur (umeshu)",
    desc: "Japanese-style plum liqueur; sweet fruit aroma, guest-friendly.",
    price: 89000,
    image: plumLiqueur,
    caffeine: "-",
    temp: "Cold",
    sweetness: "4",
    kind: "Fruit liqueur",
    isAlcoholic: true,
    baseSpirit: "Plum liqueur",
  },
];

const placeholders = [
  "https://images.unsplash.com/photo-1447933601403-0c6688de566e?auto=format&fit=crop&w=600&q=80",
  "https://images.unsplash.com/photo-1497534446932-c925b458314e?auto=format&fit=crop&w=600&q=80",
  "https://images.unsplash.com/photo-1498804103079-a6351b050096?auto=format&fit=crop&w=600&q=80",
];

const { theme } = useThemeState();
const loading = ref(true);
const error = ref("");
const toast = ref("");
const dayDrinks = ref(morningMenu);
const nightDrinks = ref(nightSpecials);

const menuMode = computed(() => (theme?.value === "night" ? "night" : "day"));
const currentMenu = computed(() =>
  menuMode.value === "night" ? nightDrinks.value : dayDrinks.value,
);

const normalizeDrink = (drink, fallbackKind = "Coffee") => ({
  ...drink,
  kind: drink.kind || fallbackKind,
  isAlcoholic:
    typeof drink.isAlcoholic === "boolean"
      ? drink.isAlcoholic
      : fallbackKind !== "Coffee",
});

const resolveAvailability = (drink) => {
  if (typeof drink.available === "boolean") {
    return drink.available;
  }

  if (typeof drink.stock === "number") {
    return drink.stock > 0;
  }

  return null;
};

const filterByTag = (items, tag) =>
  items.filter((d) =>
    (d.tags || d.Tags || []).some((t) => t?.toLowerCase() === tag),
  );

const fetchDrinks = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await getDrinks();
    const list = Array.isArray(res) && res.length ? res : [];
    const dayList = list.length ? filterByTag(list, "day") : [];
    const nightList = list.length ? filterByTag(list, "night") : [];

    dayDrinks.value = (dayList.length ? dayList : morningMenu).map((d) =>
      normalizeDrink(d, d.tags?.includes("tea") ? "Tea" : "Coffee"),
    );
    nightDrinks.value = (nightList.length ? nightList : nightSpecials).map(
      (d) => normalizeDrink(d, "Cocktail"),
    );
  } catch (err) {
    error.value =
      err?.message || "Could not load menu right now. Showing fallback menu.";
    dayDrinks.value = morningMenu.map((d) =>
      normalizeDrink(d, d.tags?.includes("tea") ? "Tea" : "Coffee"),
    );
    nightDrinks.value = nightSpecials.map((d) => normalizeDrink(d, "Cocktail"));
  } finally {
    loading.value = false;
  }
};

const imgFor = (drink, idx) =>
  drink.image || placeholders[idx % placeholders.length];
const formatCurrency = (value) =>
  `${(Number(value) || 0).toLocaleString("en-US")} VND`;
const formatSweetness = (value) => {
  if (value === undefined || value === null || value === "") return "";
  return typeof value === "number"
    ? `Sweetness ${value}/5`
    : `Sweetness ${value}`;
};
const displayCaffeine = (value) => {
  if (!value && value !== 0) return "";
  const text = String(value).trim();
  if (!text || text === "-" || text === "—") return "";
  const normalized = text.replace(/[-—\s()]/g, "").toLowerCase();
  if (!normalized || normalized.includes("teabased")) return "";
  return text;
};

// booking add-ons helper
const { add } = useCart();
const addToBookingAddons = async (drinkId) => {
  try {
    await add(drinkId, 1);
    toast.value = "Drink added to booking";
    setTimeout(() => {
      toast.value = "";
    }, 1800);
  } catch (err) {
    toast.value = "Please sign in to add drinks to your booking";
    setTimeout(() => {
      toast.value = "";
    }, 1800);
    console.warn("Add to booking failed", err);
  }
};

onMounted(fetchDrinks);
</script>

<template>
  <section class="menu-page" :class="menuMode">
    <div class="panel heading">
      <div class="hero-copy">
        <p class="kicker" :class="{ night: menuMode === 'night' }">
          {{
            menuMode === "night"
              ? "Nightfall mixology"
              : "The Creative Coffee Menu"
          }}
        </p>
        <h1 class="hero-title">
          {{
            menuMode === "night"
              ? "Mocktails, cocktails, and spirited pours."
              : "Thoughtfully brewed for every mood."
          }}
        </h1>
        <p class="lede">
          {{
            menuMode === "night"
              ? "After 6PM: low-ABV mocktails, signature cocktails, and coffee cocktails."
              : "Every cup from our Leblanc coffee lab balances smooth texture, restorative botanicals, and mindful service."
          }}
        </p>
        <div class="chips">
          <span class="chip">{{
            menuMode === "night" ? "Bar menu" : "Coffee & non-alcohol"
          }}</span>
          <span class="chip" v-if="menuMode === 'day'">Leblanc coffee lab</span>
          <span class="chip" v-else>Serving after 6PM</span>
        </div>
      </div>
      <div class="mode-banner hero-card" :class="menuMode">
        <div class="mode-label">
          <span class="mode-pill">{{
            menuMode === "night" ? "Night mode" : "Day mode"
          }}</span>
          <span class="mode-sub">
            {{
              menuMode === "night"
                ? "Mocktail & cocktail picks"
                : "Coffee, tea, juice, smoothie"
            }}
          </span>
        </div>
        <small
          >Use the Day/Night toggle in the top-right header to switch
          menus.</small
        >
      </div>
    </div>

    <div v-if="loading && menuMode === 'day'" class="state">
      Loading menu...
    </div>
    <div v-if="toast" class="state toast">{{ toast }}</div>
    <div
      v-else-if="menuMode === 'day' && !currentMenu.length"
      class="state error"
    >
      {{ error || "No menu data." }}
    </div>
    <div v-else>
      <div v-if="menuMode === 'day' && error" class="state error">
        {{ error }}
      </div>
      <div class="grid" :class="menuMode">
        <article
          v-for="(drink, idx) in currentMenu"
          :key="drink._id || drink.name"
          class="card"
          :class="menuMode"
        >
          <div class="card-media">
            <img :src="imgFor(drink, idx)" :alt="drink.name" loading="lazy" />
            <span class="kind-pill">{{
              drink.kind || (menuMode === "night" ? "Night" : "Day")
            }}</span>
          </div>
          <div class="card-body">
            <div class="card-top">
              <h3>{{ drink.name }}</h3>
              <span class="price">{{ formatCurrency(drink.price) }}</span>
            </div>
            <p class="desc">{{ drink.desc }}</p>
            <div class="meta-row">
              <span
                v-if="resolveAvailability(drink) === true"
                class="chip tone"
              >
                In stock
              </span>
              <span
                v-else-if="resolveAvailability(drink) === false"
                class="chip tone unavailable"
              >
                Sold out
              </span>
              <span v-if="displayCaffeine(drink.caffeine)" class="chip tone">{{
                displayCaffeine(drink.caffeine)
              }}</span>
              <span class="chip tone">{{ drink.temp || "-" }}</span>
              <span v-if="formatSweetness(drink.sweetness)" class="chip tone">{{
                formatSweetness(drink.sweetness)
              }}</span>
              <button
                class="btn"
                :disabled="resolveAvailability(drink) === false"
                @click="addToBookingAddons(drink._id)"
              >
                Add to booking
              </button>
            </div>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
.menu-page {
  display: flex;
  flex-direction: column;
  gap: 32px;
  color: var(--ink);
  background: var(--paper);
  padding: 32px;
}

.menu-page.night {
  color: #f5f1e8;
}

.panel {
  padding: clamp(20px, 3vw, 48px) 0;
}

.heading {
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) minmax(280px, 0.9fr);
  align-items: center;
  gap: 24px 32px;
}

.kicker {
  text-transform: uppercase;
  letter-spacing: 0.18em;
  font-size: 1.05rem;
  color: #5b4635;
  font-family: "Georgia", "Times New Roman", serif;
  font-style: italic;
  font-weight: 700;
}

:global(.theme-night) .kicker {
  color: #d6a35a !important;
}

.hero-copy {
  display: flex;
  flex-direction: column;
  gap: 14px;
  max-width: 980px;
}

.hero-title {
  font-size: clamp(2.6rem, 3vw + 1.4rem, 3.8rem);
  line-height: 1.08;
  margin: 0;
}

.lede {
  margin: 0;
  font-size: 1.08rem;
  max-width: 820px;
  opacity: 0.9;
}

.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(0, 0, 0, 0.04);
  font-weight: 750;
  letter-spacing: 0.02em;
}

.menu-page.night .chip {
  border-color: rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.08);
  color: #f5f1e8;
}

.mode-banner {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-start;
  padding: 18px 20px;
  border: 1px solid var(--cream-strong);
  border-radius: 14px;
  background: linear-gradient(
    140deg,
    rgba(235, 223, 208, 0.9),
    rgba(246, 239, 230, 0.9)
  );
  box-shadow: 0 16px 32px rgba(0, 0, 0, 0.06);
}

.mode-banner.night {
  border-color: #243145;
  background: linear-gradient(140deg, #141c2d, #0f1424);
  color: #f5f1e8;
}

.hero-card {
  min-width: 280px;
  max-width: 420px;
  justify-self: end;
}

.mode-label {
  display: grid;
  gap: 4px;
}

.mode-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 6px 12px;
  border-radius: 999px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  background: #fff;
  font-weight: 800;
  letter-spacing: 0.02em;
}

.mode-banner.night .mode-pill {
  border-color: rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.06);
  color: #f5f1e8;
}

.mode-sub {
  font-weight: 700;
  opacity: 0.88;
}

.mode-banner small {
  opacity: 0.78;
  font-size: 0.95rem;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 28px;
}

.card {
  border: 1px solid #e7ddcf;
  background: #fdfaf5;
  display: flex;
  flex-direction: column;
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.08);
  border-radius: 18px;
  overflow: hidden;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.card.night {
  background: #111827;
  border-color: #233045;
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.45);
}

.card:hover {
  transform: translateY(-6px);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.18);
}

.card-media {
  position: relative;
  height: 240px;
  overflow: hidden;
}

.card img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.25s ease;
}

.card:hover img {
  transform: scale(1.03);
}

.kind-pill {
  position: absolute;
  left: 12px;
  bottom: 12px;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(0, 0, 0, 0.7);
  color: #fff;
  font-weight: 800;
  font-size: 0.85rem;
  letter-spacing: 0.03em;
  backdrop-filter: blur(4px);
}

.card-body {
  padding: 22px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  color: var(--ink);
}

.card.night .card-top {
  color: #f5f1e8;
}

.card-top h3 {
  margin: 0;
  font-size: 1.3rem;
}

.price {
  font-weight: 900;
  letter-spacing: 0.02em;
}

.desc {
  margin: 0;
  opacity: 0.9;
}

.card.night .desc {
  color: #e7e3d9;
}

.meta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.chip.tone {
  background: rgba(0, 0, 0, 0.04);
  border-color: rgba(0, 0, 0, 0.08);
}

.chip.tone.unavailable {
  background: rgba(181, 43, 43, 0.08);
  border-color: rgba(181, 43, 43, 0.2);
  color: #8d1f1f;
}

.card.night .chip.tone {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.18);
  color: #f5f1e8;
}

.state {
  padding: 36px;
  text-align: center;
  background: var(--paper);
  border-radius: 12px;
  border: 1px solid var(--cream-strong);
}

.state.error {
  color: var(--orange-strong);
  background: #ffe8dd;
}

.state.toast {
  color: #18592a;
  background: #e9f7ec;
  border-color: #a8d7b2;
}

@media (max-width: 900px) {
  .heading {
    grid-template-columns: 1fr;
  }

  .mode-banner {
    width: 100%;
  }
}
</style>
