from unpack_entities import unpack_entities, InvalidEntitiesDetected
import unittest
import os
import shutil
import tempfile

CURRENT_DIR = os.path.dirname(os.path.abspath(__file__))


class TestUnpack(unittest.TestCase):
    def setUp(self):
        self.test_temp_dir = tempfile.mkdtemp()

    def tearDown(self):
        shutil.rmtree(self.test_temp_dir)

    def fixture_dir(self, name):
        return os.path.join(
            CURRENT_DIR,
            'fixtures', '%s_entity_packages' % name
        )

    def test_entity_package_missing_files(self):
        """Tests when the entity package is missing files"""
        with self.assertRaises(InvalidEntitiesDetected):
            unpack_entities(self.fixture_dir('bad1'), self.test_temp_dir)

    def test_entity_package_node_not_registered(self):
        """Tests when the entity package node is not registered properly"""
        with self.assertRaises(InvalidEntitiesDetected):
            unpack_entities(self.fixture_dir('bad2'), self.test_temp_dir)

    def test_succeeds(self):
        unpack_entities(self.fixture_dir('good'), self.test_temp_dir)


if __name__ == "__main__":
    unittest.main()
