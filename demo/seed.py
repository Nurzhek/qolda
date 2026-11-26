"""Fill a running qolda with demo students and listings through the real API.

    python demo/seed.py                      # app on http://localhost:8090, AUTH_DEV_LOGIN=true
    python demo/seed.py --url http://host:port

Safe to re-run: a student who already has listings is not duplicated, but their listings that have
no photo yet get one if demo/photos/<photo>.jpg (or .png/.webp) exists; `photo` is the key below.
Demo photos are from Wikimedia Commons, see demo/photos/CREDITS.md.
Afterwards post dates are spread over the last three weeks (needs `docker compose` and the db service).
"""

import argparse
import json
import mimetypes
import subprocess
import uuid
from datetime import date, timedelta
from pathlib import Path
from urllib.error import HTTPError
from urllib.request import Request, urlopen

HERE = Path(__file__).parent
TODAY = date.today()


def days(n: int) -> str:
    return (TODAY + timedelta(days=n)).isoformat()


def rent(title, category, price, deposit, until, description, photo=None, status="active"):
    return {"type": "rent", "title": title, "category": category, "price": price, "deposit": deposit,
            "available_from": days(0), "available_to": days(until), "description": description,
            "photo": photo, "status": status}


def sale(title, category, price, description, photo=None, status="active"):
    return {"type": "sale", "title": title, "category": category, "price": price, "description": description,
            "photo": photo, "status": status}


def swap(title, category, wanted, description, photo=None, status="active"):
    return {"type": "exchange", "title": title, "category": category, "wanted": wanted,
            "description": description, "photo": photo, "status": status}


# (name, e-mail local part, campus, telegram) -> listings
STUDENTS = [
    (("Aruzhan Sadykova", "aruzhan.sadykova", "dormitory", "qolda_demo_aruzhan"), [
        rent("Epson projector, 3200 lumens, HDMI", "electronics", 2500, 30000, 60,
             "Bright enough for a lit seminar room. HDMI cable and remote included, bring your own adapter for USB-C.",
             "projector"),
        sale("Calculus: Early Transcendentals, 8th ed.", "books", 7000,
             "Stewart. A few pencil notes in chapters 3-5, otherwise clean. Used for MATH 161/162.", "stewart"),
        swap("Desk lamp, warm light", "furniture", "a kettle or a small fan",
             "Warm light, two brightness levels. Moving out of the dorm and the new room has a ceiling light.", "lamp"),
    ]),
    (("Daniyar Tulegenov", "daniyar.tulegenov", "dormitory", "qolda_demo_daniyar"), [
        rent("3-person camping tent with tarp", "sports_outdoor", 3000, 25000, 75,
             "Coleman, set up in 10 minutes. Used for two Burabay trips, no holes. Pegs and repair kit inside.",
             "tent"),
        rent("Sleeping bag, comfort to -5°C", "sports_outdoor", 1500, 12000, 75,
             "Washed after every trip. Comes with a compression sack.", "sleeping-bag"),
        rent("Thermos 1.2 L and two camping mugs", "travel", 500, 4000, 75,
             "Keeps tea hot for a day. Pick up together with the tent if you want a full set.", "thermos"),
    ]),
    (("Madina Ospanova", "madina.ospanova", "family_housing", "qolda_demo_madina"), [
        rent("Casio fx-991EX scientific calculator", "study", 800, 8000, 40,
             "Allowed in most SEDS and SSH exams. Fresh battery.", "calculator"),
        sale("Lab coat, size S, barely used", "study", 3500,
             "For CHEM 101 labs. Washed, no stains. Safety goggles go with it for free.", "lab-coat"),
        sale("Organic Chemistry, Clayden, 2nd ed.", "books", 12000,
             "The good one. Hardcover, corners a bit worn, no highlights.", "chemistry-book"),
    ]),
    (("Timur Akhmetov", "timur.akhmetov", "off_campus", "qolda_demo_timur"), [
        rent("Hard-shell suitcase 70 L", "travel", 1200, 15000, 90,
             "Four wheels, TSA lock. Good for the winter break flight home.", "suitcase"),
        rent("Canon EOS 250D with 18-55 kit lens", "electronics", 5000, 80000, 45,
             "For club events or a project video. Two batteries, 64 GB card, charger.", "camera"),
        swap("Logitech G102 mouse", "electronics", "a wireless mouse",
             "Works perfectly, I just want to get rid of the cable.", "mouse"),
    ]),
    (("Aigerim Bekova", "aigerim.bekova", "dormitory", "qolda_demo_aigerim"), [
        sale("IKEA Markus office chair", "furniture", 25000,
             "Graduating in December. Mesh back, adjustable height, no squeaks. Pick up from block 22.", "chair"),
        sale("Mini fridge, 50 L", "kitchen", 30000,
             "Quiet, fits under the dorm desk. Defrosted and cleaned.", "fridge"),
        sale("Bedside rug, 80×150", "furniture", 4000, "Grey, machine washable.", "rug"),
        swap("Philips hair dryer", "other", "a straightener",
             "2000 W, cold air button works.", "hairdryer"),
    ]),
    (("Yerlan Zhumabek", "yerlan.zhumabek", "academic_block", "qolda_demo_yerlan"), [
        rent("Arduino Uno starter kit", "electronics", 1000, 10000, 50,
             "Breadboard, jumper wires, sensors, servo, LCD. Everything sorted in a box. Good for ROBT and ELCE projects.",
             "arduino"),
        rent("Soldering station Hakko FX-888D", "electronics", 1500, 25000, 50,
             "With stand, tip cleaner and lead-free solder. Return it cold, please.", "soldering"),
        sale("Introduction to Algorithms (CLRS), 3rd ed.", "books", 15000,
             "Paperback, international edition. Some highlights in the graph chapters.", "clrs"),
    ]),
    (("Sofia Kim", "sofia.kim", "dormitory", "qolda_demo_sofia"), [
        rent("Yamaha acoustic guitar F310", "music", 1500, 20000, 60,
             "New strings last month. Soft case and a capo included.", "guitar"),
        swap("Ukulele, soprano", "music", "a keyboard stand or a music stand",
             "Started guitar, so the ukulele needs a new owner.", "ukulele"),
        sale("Korean language textbooks, Sogang 1A-1B", "books", 9000,
             "Both student books and workbooks, workbooks are filled in pencil and erased.", "korean-books"),
    ]),
    (("Nursultan Omarov", "nursultan.omarov", "off_campus", "qolda_demo_nursultan"), [
        rent("Snowboard 154 cm with bindings", "sports_outdoor", 4000, 40000, 100,
             "Burton, all-mountain. Boots size 43 available too, ask in Telegram.", "snowboard"),
        rent("Ice skates, size 42", "sports_outdoor", 800, 6000, 100,
             "Sharpened in November. For the rink near Khan Shatyr.", "skates"),
        sale("Winter jacket, The North Face, size M", "clothing", 35000,
             "Down, -25°C. Bought last year, too big for me now.", "jacket"),
    ]),
    (("Kamila Iskakova", "kamila.iskakova", "dormitory", "qolda_demo_kamila"), [
        rent("Rice cooker 1.8 L", "kitchen", 700, 6000, 45,
             "For dorm dinners and club events. Steamer tray included.", "rice-cooker"),
        swap("Set of 4 plates and bowls", "kitchen", "a frying pan",
             "White ceramic, no chips.", "plates"),
        sale("Electric kettle Xiaomi", "kitchen", 6000,
             "1.7 L, keeps temperature. Original box.", "kettle"),
    ]),
    (("Alikhan Serikbay", "alikhan.serikbay", "family_housing", "qolda_demo_alikhan"), [
        rent("Bosch cordless drill", "other", 1500, 20000, 30,
             "Two batteries, set of bits. For shelves and curtains in family housing.", "drill"),
        rent("Folding table and 4 chairs", "furniture", 2000, 15000, 30,
             "For a birthday or a club stand at the activity fair.", "folding-table"),
        sale("Kids' bike, 16 inch", "sports_outdoor", 18000,
             "Training wheels included. Our son grew out of it.", "bike"),
    ]),
    (("Dilnaz Abenova", "dilnaz.abenova", "dormitory", "qolda_demo_dilnaz"), [
        rent("Graphic tablet Wacom Intuos S", "electronics", 1200, 15000, 40,
             "For design and architecture assignments. Pen with spare nibs.", "tablet"),
        sale("IELTS Cambridge 15-18 set", "books", 8000,
             "Four books with audio links. Clean, answers not written in.", "ielts-books"),
        swap("Yoga mat 6 mm", "sports_outdoor", "resistance bands",
             "Non-slip, rolled with a strap. Used for one semester of PE.", "yoga-mat"),
    ]),
    (("Arman Kassymov", "arman.kassymov", "off_campus", "qolda_demo_arman"), [
        rent("PlayStation 5 with two controllers", "electronics", 4000, 100000, 21,
             "For a weekend or a dorm tournament. FIFA 24 and Mortal Kombat 1 on the console.", "console"),
        sale("27-inch monitor Dell P2722H", "electronics", 55000,
             "IPS, 1080p, height-adjustable stand. HDMI and DisplayPort cables.", "monitor"),
        sale("Old Statistics textbook", "books", 3000, "Moved to a newer edition.", "stats-book", status="archived"),
        rent("Bluetooth speaker JBL Charge 5", "music", 1500, 20000, 21,
             "Loud enough for a club party, charges a phone too.", "speaker", status="paused"),
    ]),
]


def call(url, method="GET", body=None, token=None, raw=None, content_type="application/json"):
    data = raw if raw is not None else (json.dumps(body).encode() if body is not None else None)
    req = Request(url, data=data, method=method)
    if data is not None:
        req.add_header("Content-Type", content_type)
    if token:
        req.add_header("Authorization", f"Bearer {token}")
    try:
        with urlopen(req) as rep:
            text = rep.read().decode()
            return json.loads(text) if text.strip() else {}
    except HTTPError as e:
        raise SystemExit(f"{method} {url}: {e.code} {e.read().decode()}") from None


def upload(api, token, path: Path) -> str:
    boundary = uuid.uuid4().hex
    mime = mimetypes.guess_type(path.name)[0] or "application/octet-stream"
    body = (f"--{boundary}\r\nContent-Disposition: form-data; name=\"file\"; filename=\"{path.name}\"\r\n"
            f"Content-Type: {mime}\r\n\r\n").encode() + path.read_bytes() + f"\r\n--{boundary}--\r\n".encode()
    return call(f"{api}/file", "POST", token=token, raw=body,
                content_type=f"multipart/form-data; boundary={boundary}")["url"]


def add_missing_photos(api, token, mine, listings) -> int:
    by_title = {item["title"]: item for item in listings}
    added = 0
    for listing in mine:
        item = by_title.get(listing["title"])
        path = find_photo(item["photo"]) if item and not listing["photos"] else None
        if path:
            call(f"{api}/listing/{listing['id']}", "PUT", {"photos": [upload(api, token, path)]}, token=token)
            added += 1
    return added


def find_photo(name):
    if not name:
        return None
    return next((p for ext in ("jpg", "jpeg", "png", "webp") if (p := HERE / "photos" / f"{name}.{ext}").exists()), None)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--url", default="http://localhost:8090")
    ap.add_argument("--domain", default="nu.edu.kz")
    ap.add_argument("--no-backdate", action="store_true")
    args = ap.parse_args()
    api = args.url.rstrip("/") + "/api"

    if not call(f"{api}/dict").get("dev_login"):
        raise SystemExit("dev sign-in is off: set AUTH_DEV_LOGIN=true for the app")

    created, photos = [], 0
    for (name, local, campus, telegram), listings in STUDENTS:
        token = call(f"{api}/auth/dev", "POST", {"email": f"{local}@{args.domain}", "name": name})["access_token"]
        mine = call(f"{api}/listing?mine=true&list_params.page_size=100", token=token)["results"]
        if mine:
            added = add_missing_photos(api, token, mine, listings)
            photos += added
            print(f"{name}: already seeded, {added} photos added")
            continue
        call(f"{api}/me", "PUT", {"name": name, "campus": campus, "telegram": telegram}, token=token)

        for item in listings:
            body = {k: v for k, v in item.items() if k not in ("photo", "status")}
            body["campus"] = campus
            path = find_photo(item["photo"])
            body["photos"] = [upload(api, token, path)] if path else []
            photos += bool(path)
            listing_id = call(f"{api}/listing", "POST", body, token=token)["id"]
            if item["status"] != "active":
                call(f"{api}/listing/{listing_id}/status", "PUT", {"status": item["status"]}, token=token)
            created.append(listing_id)
        print(f"{name}: {len(listings)} listings")

    print(f"created {len(created)} listings, {photos} with photos")
    if created and not args.no_backdate:
        backdate(created)


def backdate(ids):
    """Spread post dates over three weeks so the catalogue doesn't say "this minute" everywhere."""
    id_list = ",".join(str(i) for i in ids)
    sql = (f"update listing set created_at = d, modified_at = d from (select id as lid, "
           f"now() - random() * interval '21 days' as d from listing where id in ({id_list})) x where id = x.lid; "
           f"update usr set created_at = now() - (interval '30 days' + random() * interval '300 days') "
           f"where id in (select owner_id from listing where id in ({id_list}));")
    cmd = ["docker", "compose", "exec", "-T", "db", "psql", "-U", "qolda", "-q", "-c", sql]
    try:
        subprocess.run(cmd, check=True, cwd=HERE.parent, capture_output=True)
        print("post dates spread over the last three weeks")
    except (OSError, subprocess.CalledProcessError) as e:
        print(f"could not backdate (listings are fine, just all posted now): {e}")


if __name__ == "__main__":
    main()
