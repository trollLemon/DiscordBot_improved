import requests
import unittest



class EndPointIntegrationTest(unittest.TestCase):

    def test_invert(self):
        api_call_good = "http://127.0.0.1:8000/api/invertedImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D"
        api_call_bad_url =  "http://127.0.0.1:8000/api/invertedImage/notanimagelink/"

        r = requests.get(url=api_call_good)
        self.assertEqual(r.status_code, 200)
        r = requests.get(url=api_call_bad_url)
        self.assertEqual(r.status_code, 422)

    def test_saturate(self):
        api_call_good = "http://127.0.0.1:8000/api/saturatedImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/10"
        api_call_bad_url =  "http://127.0.0.1:8000/api/saturatedImage/notanimagelink/3"

        r = requests.get(url=api_call_good)
        self.assertEqual(r.status_code, 200)
        r = requests.get(url=api_call_bad_url)
        self.assertEqual(r.status_code, 422)

    def test_edge_detect(self):
        api_call_good = "http://127.0.0.1:8000/api/edgeImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/100/200"
        api_call_bad_param =  "http://127.0.0.1:8000/api/edgeImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/-1/3"
        api_call_bad_url =  "http://127.0.0.1:8000/api/edgeImage/notanimagelink/3/3"

        r = requests.get(url=api_call_good)
        self.assertEqual(r.status_code, 200)
        r = requests.get(url=api_call_bad_param)
        self.assertEqual(r.status_code, 400)
        r = requests.get(url=api_call_bad_url)
        self.assertEqual(r.status_code, 422)

    def test_dilate(self):
        api_call_good = "http://127.0.0.1:8000/api/dilatedImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/3/3"
        api_call_bad_param =  "http://127.0.0.1:8000/api/dilatedImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/0/3"
        api_call_bad_url =  "http://127.0.0.1:8000/api/dilatedImage/notanimagelink/3/3"

        r = requests.get(url=api_call_good)
        self.assertEqual(r.status_code, 200)
        r = requests.get(url=api_call_bad_param)
        self.assertEqual(r.status_code, 400)
        r = requests.get(url=api_call_bad_url)
        self.assertEqual(r.status_code, 422)

    def test_erode(self):
        api_call_good = "http://127.0.0.1:8000/api/erodedImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/3/3"
        api_call_bad_param =  "http://127.0.0.1:8000/api/erodedImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/0/3"
        api_call_bad_url =  "http://127.0.0.1:8000/api/erodedImage/notanimagelink/3/3"

        r = requests.get(url=api_call_good)
        self.assertEqual(r.status_code, 200)
        r = requests.get(url=api_call_bad_param)
        self.assertEqual(r.status_code, 400)
        r = requests.get(url=api_call_bad_url)
        self.assertEqual(r.status_code, 422)

    def test_reduce(self):

        api_call_good = "http://127.0.0.1:8000/api/reducedImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/0.1"
        api_call_bad_param =  "http://127.0.0.1:8000/api/reducedImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/6"
        api_call_bad_url =  "http://127.0.0.1:8000/api/reducedImage/notanimagelink/0.1"


        r = requests.get(url=api_call_good)
        self.assertEqual(r.status_code, 200)
        r = requests.get(url=api_call_bad_param)
        self.assertEqual(r.status_code, 400)
        r = requests.get(url=api_call_bad_url)
        self.assertEqual(r.status_code, 422)

    def test_random_kernel(self):


        api_call_good = "http://127.0.0.1:8000/api/randomFilteredImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/4/0/6/raw"
        api_call_bad_param =  "http://127.0.0.1:8000/api/randomFilteredImage/https%3A%2F%2Fmedia.istockphoto.com%2Fid%2F1648451466%2Fphoto%2Fmale-medical-worker-pipetting-chemical-in-test-tube-while-working-in-laboratory.jpg%3Fs%3D612x612%26w%3Dis%26k%3D20%26c%3DyBdsqEM7Lce__SWpRdqyljvb9uxhNho2K0FKohWJdzo%3D/-2/0/6/norm"
        api_call_bad_url =  "http://127.0.0.1:8000/api/randomFilteredImage/notanimagelink/3/0/4/raw"


        r = requests.get(url=api_call_good)
        self.assertEqual(r.status_code, 200)
        r = requests.get(url=api_call_bad_param)
        self.assertEqual(r.status_code, 400)
        r = requests.get(url=api_call_bad_url)
        self.assertEqual(r.status_code, 422)


    def test_add_text(self):
        pass


if __name__ == '__main__':
    unittest.main()
