const PATH = window.location.pathname;
const boot = window.GOPHER_BOOT || {spawn: false, raceActive: false, timeLeft: 0};
let raceActive = boot.raceActive;
let timeLeft = boot.timeLeft;

window.addEventListener("error", e => logError("window.onerror", e.error || e));
window.addEventListener("unhandledrejection",
    e => logError("unhandled‑promise", e.reason));

const randVH = () => Math.floor(Math.random() * 81) + "vh";
const randVW = () => Math.floor(Math.random() * 81) + "vw";

function logError(context, err) {
    console.error(`[Gopher][${context}]`, err);

    const box = document.getElementById("gopher-error");
    if (box) {
        box.textContent = `⚠ ${context}: ${err.message || err}`;
        box.classList.remove("hidden");
    }
}

async function safeFetch(url, options = {}, context = "fetch") {
    try {
        const res = await fetch(url, options);
        if (!res.ok) {
            throw new Error(`${url} → ${res.status} ${res.statusText}`);
        }
        return res;
    } catch (err) {
        logError(context, err);
        throw err;
    }
}

function wire(img) {
    img.onclick = async () => {
        img.remove();
        try {
            const res = await safeFetch("/click", {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({path: PATH}),
                keepalive: true
            }, "click");
            const data = await res.json();
            // Update Score Banner
            setScore("#score-display",  data.score);
            setScore("#total-count",   data.totalCatches);
            function setScore(id, val){
                const el = document.querySelector(id);
                if (el){
                    el.textContent = val;
                    el.classList.add("bump");
                    setTimeout(()=>el.classList.remove("bump"), 350);
                }
            }
            // Display Fact Box
            const fact = document.getElementById("fact-box");
            if (fact && !raceActive) {
                fact.textContent = data.funFact;
                fact.classList.remove("hidden");
                setTimeout(()=> fact.classList.add("hidden"), 5000);
            }
            // Update Race Score
            if (raceActive) {
                setTextById("race-score", `${data.raceScore}`);
            }

            if (data.spawnAgain) spawnGopher();
        } catch (_) {
        }
    };
}

/* ---------- Gopher element ---------- */
function spawnGopher() {
    const g = document.createElement("img");
    g.src = "/static/img/gopher.png";
    g.className = "gopher";
    if (raceActive) g.classList.add("moving-gopher");

    g.style.top = randVH();
    g.style.left = randVW();
    document.body.appendChild(g);
    wire(g);
    return g;
}

async function startRaceBtn() {
    if (raceActive) return;
    try {
        await safeFetch("/start-race", {method: "POST", keepalive: true});
        raceActive = true;
        timeLeft = 30;
        setTextById("race-score", "0");
        document.getElementById("race-score").classList.remove("hidden");
        spawnGopher();
        startCountdown();
    } catch (_) {
    }
}

function updateArc() {
    const pct = timeLeft / 30;
    const progress = document.querySelector(".progress");
    const FULL = 300;
    progress.style.strokeDashoffset = (FULL * (1 - pct)).toString();

    if (timeLeft <= 5)      progress.style.stroke = "#dc2626";
    else if (timeLeft <=15) progress.style.stroke = "#f97316";
    else                    progress.style.stroke = "#22c55e";
}

function startCountdown() {
    const clock = document.getElementById("race-countdown")
    const digits = document.getElementById("race-timer");

    // First Frame
    digits.textContent = timeLeft;
    clock.style.display = "revert"
    updateArc()
    const iv = setInterval(() => {
        timeLeft--;
        digits.textContent = timeLeft;

        // Timer Pulse
        digits.classList.add("pulse");
        setTimeout(() => digits.classList.remove("pulse"), 150);

        // Shrink Arc with Color
        updateArc()

        // End
        if (timeLeft <= 0) {
            clearInterval(iv);
            digits.textContent = "0"
            clock.style.display = "none"
            raceActive = false;
            endRace();
        }
    }, 1000);
}

function setTextById(id, val) {
    const el = document.getElementById(id);
    if (el) el.textContent = val;
}

async function endRace() {
    try {
        const res = await safeFetch("/end-race",
            {method: "POST", keepalive: true}, "end-race");
        const d = await res.json();

        document.getElementById("race-score").classList.add("hidden");
        setTextById("race-score", "0");
        setTextById("session-race", d.sessionHigh);
        setTextById("global-race", d.globalHigh);
        setTextById("total-count", `${d.totalCatches}`);
        showToast(`Race over, you scored ${d.raceScore}!`);
    } catch (err) {
        logError("end-race", err);
        showToast("⚠ Race finished, but saving stats failed.");
    }
}

function showToast(msg, ms = 3000) {
    let toast = document.getElementById("gopher-toast");
    if (!toast) {
        toast = document.createElement("div");
        toast.id = "gopher-toast";
        toast.className = "toast hidden";
        toast.innerHTML = "<span class='toast-msg'></span>";
        document.body.appendChild(toast);
    }
    toast.querySelector(".toast-msg").textContent = msg;
    toast.style.display = "flex";
    toast.querySelectorAll(".firework").forEach(fw => fw.innerHTML = "");
    toast.querySelectorAll(".firework").forEach(fw => {
        for (let i = 0; i < 8; i++) {
            const ray = document.createElement("span");
            ray.className = "ray";
            ray.style.transform = `translate(-50%,-50%) rotate(${i * 45}deg)`;
            const dot = document.createElement("span");
            dot.className = "dot";
            ray.appendChild(dot);
            fw.appendChild(ray);
        }
    });
    setTimeout(() => toast.style.display = "none", ms);
}

/* ---------- boot ---------- */
document.addEventListener("DOMContentLoaded", () => {
    const bootImg = document.getElementById("boot-gopher");
    if (bootImg) wire(bootImg);
    if (raceActive && timeLeft > 0) {
        startCountdown();
    }

    const raceBtn = document.getElementById("start-race-btn");
    if (raceBtn) raceBtn.addEventListener("click", startRaceBtn);
});
