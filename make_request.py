import os
import requests

URL = "http://localhost:8080/generate_text"
API_KEY = "TODO Later"


def generate_text(
    image_path: str, text: str, size: int, x: int, y: int, r: int, g: int, b: int
):
    file_ext = "png" if image_path.lower().endswith(".png") else "jpg"
    with open(image_path, "rb") as image_file:
        files = {"image_data": (f"image.{file_ext}", image_file, f"image/{file_ext}")}
        data = {
            "api_key": API_KEY,
            "text": text,
            "x": x,
            "y": y,
            "r": r,
            "g": g,
            "b": b,
            "size": size,
        }
        response = requests.post(URL, files=files, data=data)

    if response.status_code == 200:
        output_filename = os.path.join(
            "generated", f"new_{os.path.basename(image_path)}"
        )
        os.makedirs(os.path.dirname(output_filename), exist_ok=True)
        with open(output_filename, "wb") as f:
            f.write(response.content)
        print(f"Image saved as {output_filename}")
    else:
        print(f"Error: {response.status_code} - {response.text}")


if __name__ == "__main__":
    image_names = [
        "blacksand.jpg",
        "bluepink.jpg",
        "green with noise.jpg",
        "m3pro.png",
        "purple wave.jpg",
        "purple.jpg",
        "whitesand.jpg",
    ]
    for name in image_names:
        generate_text(f"images/{name}", "Some text!", 150, 300, 300, 255, 255, 255)
