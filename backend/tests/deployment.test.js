const { expect } = require('chai');
const sinon = require('sinon');
const path = require('path');
const fs = require('fs');
const { 
  generateTerraformConfig, 
  generateTopology,
  ensureDirectoryExists,
  writeToFile,
  readFile
} = require('../deployment');

describe('Deployment Module', () => {
  describe('generateTerraformConfig', () => {
    it('should generate AWS terraform configuration correctly', () => {
      const config = {
        cloudProvider: 'aws',
        region: 'us-east-1',
        az: 'us-east-1a',
        vpc: {
          name: 'test-vpc',
          cidr: '10.0.0.0/16',
          enableDnsSupport: true,
          enableDnsHostnames: true
        },
        subnet: {
          name: 'test-subnet',
          cidr: '10.0.1.0/24',
          mapPublicIpOnLaunch: true
        },
        components: [
          { value: 'compute', name: '弹性计算实例' }
        ],
        componentProperties: {
          'compute': {
            'instance_count': '2',
            'instance_type': 'medium'
          }
        }
      };
      
      const terraformConfig = generateTerraformConfig(config);
      
      // 验证包含AWS提供商配置
      expect(terraformConfig).to.include('provider "aws"');
      expect(terraformConfig).to.include('region = "us-east-1"');
      
      // 验证包含VPC配置
      expect(terraformConfig).to.include('resource "aws_vpc" "test-vpc"');
      expect(terraformConfig).to.include('cidr_block           = "10.0.0.0/16"');
      
      // 验证包含子网配置
      expect(terraformConfig).to.include('resource "aws_subnet" "test-subnet"');
      expect(terraformConfig).to.include('cidr_block              = "10.0.1.0/24"');
      expect(terraformConfig).to.include('availability_zone       = "us-east-1a"');
      
      // 验证包含计算实例配置
      expect(terraformConfig).to.include('resource "aws_instance" "web"');
      expect(terraformConfig).to.include('count         = 2');
      expect(terraformConfig).to.include('instance_type = "t3.medium"');
    });
    
    it('should generate Azure terraform configuration correctly', () => {
      const config = {
        cloudProvider: 'azure',
        region: 'eastus',
        az: 'eastus-1',
        vpc: {
          name: 'test-vnet',
          cidr: '10.0.0.0/16'
        },
        subnet: {
          name: 'test-subnet',
          cidr: '10.0.1.0/24'
        },
        components: [
          { value: 'object-storage', name: '对象存储' }
        ],
        componentProperties: {
          'object-storage': {
            'bucket_name': 'teststorage'
          }
        }
      };
      
      const terraformConfig = generateTerraformConfig(config);
      
      // 验证包含Azure提供商配置
      expect(terraformConfig).to.include('provider "azurerm"');
      expect(terraformConfig).to.include('features {}');
      
      // 验证包含资源组和VNet配置
      expect(terraformConfig).to.include('resource "azurerm_resource_group"');
      expect(terraformConfig).to.include('resource "azurerm_virtual_network" "test-vnet"');
      expect(terraformConfig).to.include('address_space       = ["10.0.0.0/16"]');
      
      // 验证包含子网配置
      expect(terraformConfig).to.include('resource "azurerm_subnet" "test-subnet"');
      expect(terraformConfig).to.include('address_prefixes     = ["10.0.1.0/24"]');
      
      // 验证包含存储账户配置
      expect(terraformConfig).to.include('resource "azurerm_storage_account"');
    });
    
    it('should generate Alicloud terraform configuration correctly', () => {
      const config = {
        cloudProvider: 'alicloud',
        region: 'cn-hangzhou',
        az: 'cn-hangzhou-a',
        vpc: {
          name: 'test-vpc',
          cidr: '192.168.0.0/16'
        },
        subnet: {
          name: 'test-vswitch',
          cidr: '192.168.1.0/24'
        },
        components: [],
        componentProperties: {}
      };
      
      const terraformConfig = generateTerraformConfig(config);
      
      // 验证包含阿里云提供商配置
      expect(terraformConfig).to.include('provider "alicloud"');
      expect(terraformConfig).to.include('region = "cn-hangzhou"');
      
      // 验证包含VPC配置
      expect(terraformConfig).to.include('resource "alicloud_vpc" "test-vpc"');
      expect(terraformConfig).to.include('cidr_block = "192.168.0.0/16"');
      
      // 验证包含交换机配置
      expect(terraformConfig).to.include('resource "alicloud_vswitch" "test-vswitch"');
      expect(terraformConfig).to.include('cidr_block = "192.168.1.0/24"');
      expect(terraformConfig).to.include('zone_id    = "cn-hangzhou-a"');
    });
  });
  
  describe('generateTopology', () => {
    it('should generate topology with correct nodes and edges', () => {
      const config = {
        cloudProvider: 'aws',
        region: 'us-east-1',
        az: 'us-east-1a',
        vpc: {
          name: 'test-vpc',
          cidr: '10.0.0.0/16'
        },
        subnet: {
          name: 'test-subnet',
          cidr: '10.0.1.0/24'
        },
        components: [
          { value: 'load-balancer', name: '负载均衡器' },
          { value: 'compute', name: '弹性计算实例' }
        ],
        componentProperties: {}
      };
      
      const topology = generateTopology(config);
      
      // 验证节点数量
      expect(topology.nodes.length).to.equal(4); // cloud, vpc, subnet, 2 components
      
      // 验证包含云节点
      const cloudNode = topology.nodes.find(node => node.type === 'cloud');
      expect(cloudNode).to.exist;
      expect(cloudNode.name).to.equal('AWS');
      
      // 验证包含VPC节点
      const vpcNode = topology.nodes.find(node => node.type === 'vpc');
      expect(vpcNode).to.exist;
      expect(vpcNode.name).to.equal('test-vpc');
      expect(vpcNode.data.cidr).to.equal('10.0.0.0/16');
      
      // 验证包含子网节点
      const subnetNode = topology.nodes.find(node => node.type === 'subnet');
      expect(subnetNode).to.exist;
      expect(subnetNode.name).to.equal('test-subnet');
      expect(subnetNode.data.cidr).to.equal('10.0.1.0/24');
      
      // 验证边的数量和连接关系
      expect(topology.edges.length).to.be.at.least(3); // cloud->vpc, vpc->subnet, subnet->components
      
      // 验证云到VPC的边
      const cloudToVpcEdge = topology.edges.find(edge => 
        edge.source === 'cloud' && edge.target === 'vpc');
      expect(cloudToVpcEdge).to.exist;
      
      // 验证VPC到子网的边
      const vpcToSubnetEdge = topology.edges.find(edge => 
        edge.source === 'vpc' && edge.target === 'subnet');
      expect(vpcToSubnetEdge).to.exist;
    });
    
    it('should create connections between related components', () => {
      const config = {
        cloudProvider: 'aws',
        region: 'us-east-1',
        az: 'us-east-1a',
        vpc: {
          name: 'test-vpc',
          cidr: '10.0.0.0/16'
        },
        subnet: {
          name: 'test-subnet',
          cidr: '10.0.1.0/24'
        },
        components: [
          { value: 'load-balancer', name: '负载均衡器' },
          { value: 'compute', name: '弹性计算实例' },
          { value: 'cdn', name: 'CDN' },
          { value: 'object-storage', name: '对象存储' }
        ],
        componentProperties: {}
      };
      
      const topology = generateTopology(config);
      
      // 验证负载均衡器到计算实例的连接
      const lbToComputeEdge = topology.edges.find(edge => 
        edge.source.includes('component') && 
        edge.target.includes('component') && 
        edge.type === 'connects');
      expect(lbToComputeEdge).to.exist;
      
      // 验证CDN到对象存储的连接
      const cdnToStorageEdge = topology.edges.find(edge => 
        edge.source.includes('component') && 
        edge.target.includes('component') && 
        edge.type === 'origin');
      expect(cdnToStorageEdge).to.exist;
    });
  });
  
  describe('File operations', () => {
    let sandbox;
    
    beforeEach(() => {
      sandbox = sinon.createSandbox();
    });
    
    afterEach(() => {
      sandbox.restore();
    });
    
    it('should ensure directory exists', () => {
      const existsStub = sandbox.stub(fs, 'existsSync').returns(false);
      const mkdirStub = sandbox.stub(fs, 'mkdirSync');
      
      ensureDirectoryExists('/test/dir');
      
      expect(existsStub.calledWith('/test/dir')).to.be.true;
      expect(mkdirStub.calledWith('/test/dir', { recursive: true })).to.be.true;
    });
    
    it('should not create directory if it already exists', () => {
      const existsStub = sandbox.stub(fs, 'existsSync').returns(true);
      const mkdirStub = sandbox.stub(fs, 'mkdirSync');
      
      ensureDirectoryExists('/test/dir');
      
      expect(existsStub.calledWith('/test/dir')).to.be.true;
      expect(mkdirStub.called).to.be.false;
    });
    
    it('should write to file', async () => {
      const writeFileStub = sandbox.stub(fs, 'writeFile').callsFake((path, content, encoding, callback) => {
        callback(null);
      });
      
      await writeToFile('/test/file.txt', 'test content');
      
      expect(writeFileStub.calledWith('/test/file.txt', 'test content', 'utf8')).to.be.true;
    });
    
    it('should read from file', async () => {
      const readFileStub = sandbox.stub(fs, 'readFile').callsFake((path, encoding, callback) => {
        callback(null, 'test content');
      });
      
      const content = await readFile('/test/file.txt');
      
      expect(readFileStub.calledWith('/test/file.txt', 'utf8')).to.be.true;
      expect(content).to.equal('test content');
    });
  });
});
